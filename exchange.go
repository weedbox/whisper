package whisper

type Exchange interface {
	Init(w Whisper) error
	AssertStream(streamName string, subjects []string) error
	DeleteStream(streamName string) error
	EmitSignal(subject string, payload []byte) error
	EmitEvent(subject string, payload []byte) error
	SendBatchEvents(events []*Event) error
	ChanSubscribe(subject string, channels chan Msg, durable string) (Subscription, error)
	ChanQueueSubscribe(subject string, queue string, channels chan Msg) (Subscription, error)
	ChanWatch(subject string, channels chan Msg) (Subscription, error)
}

type Msg interface {
	Payload() []byte
	Ack() error
	Nak() error
}

type Subscription interface {
	Unsubscribe() error
}

type Event struct {
	Subject string
	Payload []byte
}
