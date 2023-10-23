package main

import "fmt"

const (
	MEMBER_AGENT_DURABLE_FORMAT = "%s_member_agent_%s" // <domain>_member_agent_<durable>
)

type MemberAgent interface {
	ChanWatch(ch chan Msg) error
	ChanSubscribe(durable string, ch chan Msg) error
}

type memberAgent struct {
	w        Whisper
	memberID string
}

func NewMemberAgent(w Whisper, memberID string) MemberAgent {
	return &memberAgent{
		w:        w,
		memberID: memberID,
	}
}

func (ma *memberAgent) ChanWatch(ch chan Msg) error {

	subject := fmt.Sprintf(MEMBER_SUBJECT_FORMAT, ma.w.Domain(), "*", ma.memberID)

	_, err := ma.w.Exchange().ChanWatch(subject, ch)
	if err != nil {
		return err
	}

	return nil
}

func (ma *memberAgent) ChanSubscribe(durable string, ch chan Msg) error {

	subject := fmt.Sprintf(MEMBER_SUBJECT_FORMAT, ma.w.Domain(), "*", ma.memberID)
	durableName := fmt.Sprintf(MEMBER_AGENT_DURABLE_FORMAT, ma.w.Domain(), durable)

	_, err := ma.w.Exchange().ChanSubscribe(subject, ch, durableName)
	if err != nil {
		return err
	}

	return nil
}
