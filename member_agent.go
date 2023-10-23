package main

import "fmt"

type MemberAgent interface {
	ChanWatch(ch chan Msg) error
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
