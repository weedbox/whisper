package whisper

import (
	"fmt"

	"github.com/lithammer/go-jump-consistent-hash"
)

const (
	MESSAGE_DISPATCHER_DURABLE_FORMAT = "%s_message_dispatcher"  // <domain>_message_dispatcher
	MEMBER_SUBJECT_FORMAT             = "%s.member.bucket.%v.%s" // <domain>.member.bucket.<hash>.<mid>
	MEMBER_STREAM_FORMAT              = "%s_MEMBER_%d_MSG"       // <domain>_MEMBER_<hash>_MSG
)

type MessageDispatcher interface {
	Init() error
	ClearStream() error
	Subscribe(buckets []int32) error
	Dispatch(memberID string, payload []byte) error
	DispatchAll(members []string, payload []byte) error
	Close()
}

type messageDispatcher struct {
	w        Whisper
	subs     map[int32]Subscription
	channels map[int32]chan Msg
}

func NewMessageDispatcher(w Whisper) MessageDispatcher {
	return &messageDispatcher{
		w:        w,
		channels: make(map[int32]chan Msg),
		subs:     make(map[int32]Subscription),
	}
}

func (md *messageDispatcher) subscribe(bucket int32) error {

	subject := fmt.Sprintf(GROUP_SUBJECT_FORMAT, md.w.Domain(), bucket, "*")
	durable := fmt.Sprintf(MESSAGE_DISPATCHER_DURABLE_FORMAT, md.w.Domain())

	ch := make(chan Msg, 2048)
	md.channels[bucket] = ch

	sub, err := md.w.Exchange().ChanSubscribe(subject, ch, durable)
	if err != nil {
		return err
	}

	md.subs[bucket] = sub

	go func() {
		for m := range ch {
			err := md.HandleMessage(m)
			if err != nil {
				m.Nak()
				continue
			}

			m.Ack()
		}
	}()

	return nil
}

func (md *messageDispatcher) Subscribe(buckets []int32) error {

	for _, b := range buckets {
		md.subscribe(b)
	}

	return nil
}

func (md *messageDispatcher) Init() error {

	// Initializing streams for members
	for i := 0; i < int(md.w.BucketSize()); i++ {

		streamName := fmt.Sprintf(MEMBER_STREAM_FORMAT, md.w.Domain(), i)
		subject := fmt.Sprintf(MEMBER_SUBJECT_FORMAT, md.w.Domain(), i, ">")

		err := md.w.Exchange().AssertStream(streamName, []string{subject})
		if err != nil {
			return err
		}
	}

	return nil
}

func (md *messageDispatcher) ClearStream() error {

	for i := 0; i < int(md.w.BucketSize()); i++ {

		streamName := fmt.Sprintf(MEMBER_STREAM_FORMAT, md.w.Domain(), i)

		err := md.w.Exchange().DeleteStream(streamName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (md *messageDispatcher) Close() {

	for _, ch := range md.channels {
		close(ch)
	}
}

func (md *messageDispatcher) HandleMessage(m Msg) error {

	fmt.Println(string(m.Payload()))

	msg, err := ParseMessage(m.Payload())
	if err != nil {

		// Ignore if invalid message
		if err == ErrInvalidMessageFormat {
			return nil
		}

		return err
	}

	// Getting all of members in the group
	rule := md.w.GroupResolver().GetGroupRule(msg.Meta.Group)

	fmt.Println(rule)

	if rule == nil {
		// No such group
		return nil
	}

	// Dispatch to all members in this group
	members := rule.GetMembers()
	return md.DispatchAll(members, m.Payload())
}

func (md *messageDispatcher) DispatchAll(members []string, payload []byte) error {

	// Preparing events
	var events []*Event
	for _, memberID := range members {

		h := jump.HashString(memberID, md.w.BucketSize(), jump.NewCRC64())

		subject := fmt.Sprintf(MEMBER_SUBJECT_FORMAT, md.w.Domain(), h, memberID)

		events = append(events, &Event{
			Subject: subject,
			Payload: payload,
		})
	}

	err := md.w.Exchange().SendBatchEvents(events)
	if err != nil {
		return err
	}

	return nil
}

func (md *messageDispatcher) Dispatch(memberID string, payload []byte) error {

	h := jump.HashString(memberID, md.w.BucketSize(), jump.NewCRC64())

	subject := fmt.Sprintf(MEMBER_SUBJECT_FORMAT, md.w.Domain(), h, memberID)
	/*
		err := md.w.Exchange().EmitEvent(subject, payload)
		if err != nil {
			return err
		}
	*/
	err := md.w.Exchange().EmitSignal(subject, payload)
	if err != nil {
		return err
	}

	return nil
}
