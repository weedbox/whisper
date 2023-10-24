package whisper

import (
	"fmt"

	"github.com/lithammer/go-jump-consistent-hash"
)

const (
	GROUP_SUBJECT_FORMAT = "%s.group.bucket.%v.%s" // <domain>.group.bucket.<hash>.<gid>
	GROUP_STREAM_FORMAT  = "%s_GROUP_%d_MSG"       // <domain>_GROUP_<hash>_MSG
)

type MessageReceiver interface {
	Init() error
	ClearStream() error
	Receive(payload []byte) error
}

type messageReceiver struct {
	w Whisper
}

func NewMessageReceiver(w Whisper) MessageReceiver {
	return &messageReceiver{
		w: w,
	}
}

func (mr *messageReceiver) Init() error {

	// Initializing streams for groups
	for i := 0; i < int(mr.w.BucketSize()); i++ {

		streamName := fmt.Sprintf(GROUP_STREAM_FORMAT, mr.w.Domain(), i)
		subject := fmt.Sprintf(GROUP_SUBJECT_FORMAT, mr.w.Domain(), i, ">")

		err := mr.w.Exchange().AssertStream(streamName, []string{subject})
		if err != nil {
			return err
		}
	}

	return nil
}

func (mr *messageReceiver) ClearStream() error {

	for i := 0; i < int(mr.w.BucketSize()); i++ {

		streamName := fmt.Sprintf(GROUP_STREAM_FORMAT, mr.w.Domain(), i)

		err := mr.w.Exchange().DeleteStream(streamName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mr *messageReceiver) Receive(payload []byte) error {

	msg, err := ParseMessage(payload)
	if err != nil {

		// Ignore if invalid message
		if err == ErrInvalidMessageFormat {
			return nil
		}

		return err
	}

	h := jump.HashString(msg.Meta.Group, mr.w.BucketSize(), jump.NewCRC64())

	subject := fmt.Sprintf(GROUP_SUBJECT_FORMAT, mr.w.Domain(), h, msg.Meta.Group)
	/*
		err := mr.w.Exchange().EmitEvent(subject, payload)
		if err != nil {
			return err
		}
	*/
	err = mr.w.Exchange().EmitSignal(subject, payload)
	if err != nil {
		return err
	}

	return nil
}
