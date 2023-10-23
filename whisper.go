package main

import (
	"errors"
)

var (
	ErrInvalidMessageFormat = errors.New("whisper: invalid message format")
)

type Whisper interface {
	Init() error
	Exchange() Exchange
	GroupManager() GroupManager
	Domain() string
	BucketSize() int32
	NewMessageReceiver() (MessageReceiver, error)
	NewMessageDispatcher() (MessageDispatcher, error)
	NewMemberAgent(memberID string) (MemberAgent, error)
}

type whisper struct {
	ex         Exchange
	gm         GroupManager
	domain     string
	bucketSize int32
}

type Opt func(*whisper)

func WithExchange(ex Exchange) Opt {
	return func(w *whisper) {
		w.ex = ex
	}
}

func WithGroupManager(gm GroupManager) Opt {
	return func(w *whisper) {
		w.gm = gm
	}
}

func WithDomain(domain string) Opt {
	return func(w *whisper) {
		w.domain = domain
	}
}

func WithBucketSize(size int32) Opt {
	return func(w *whisper) {
		w.bucketSize = size
	}
}

func NewWhisper(opts ...Opt) Whisper {

	w := &whisper{}

	for _, o := range opts {
		o(w)
	}

	return w
}

func (w *whisper) Init() error {

	err := w.ex.Init(w)
	if err != nil {
		return err
	}

	err = w.gm.Init(w)
	if err != nil {
		return err
	}

	return nil
}

func (w *whisper) Exchange() Exchange {
	return w.ex
}

func (w *whisper) GroupManager() GroupManager {
	return w.gm
}

func (w *whisper) Domain() string {
	return w.domain
}

func (w *whisper) BucketSize() int32 {
	return w.bucketSize
}

func (w *whisper) NewMessageReceiver() (MessageReceiver, error) {

	mr := NewMessageReceiver(w)

	err := mr.Init()
	if err != nil {
		return nil, err
	}

	return mr, nil
}

func (w *whisper) NewMessageDispatcher() (MessageDispatcher, error) {

	md := NewMessageDispatcher(w)

	err := md.Init()
	if err != nil {
		return nil, err
	}

	return md, nil
}

func (w *whisper) NewMemberAgent(memberID string) (MemberAgent, error) {

	ma := NewMemberAgent(w, memberID)

	return ma, nil
}