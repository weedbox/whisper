package main

import (
	"github.com/nats-io/nats.go"
)

type NATSMsg struct {
	m *nats.Msg
}

func (m *NATSMsg) Payload() []byte {
	return m.m.Data
}

func (m *NATSMsg) Ack() error {
	return m.m.Ack()
}

func (m *NATSMsg) Nak() error {
	return m.m.Nak()
}

type NATSSubscription struct {
	sub *nats.Subscription
	ch  chan *nats.Msg
}

func (s *NATSSubscription) Unsubscribe() error {
	return s.sub.Unsubscribe()
}

func (s *NATSSubscription) Close() {
	close(s.ch)
}

type NATSExchange struct {
	conn *nats.Conn
}

func NewNATSExchange(conn *nats.Conn) Exchange {
	return &NATSExchange{
		conn: conn,
	}
}

func (ex *NATSExchange) AssertStream(streamName string, subjects []string) error {

	js, _ := ex.conn.JetStream()

	_, err := js.StreamInfo(streamName)
	if err != nil && err != nats.ErrStreamNotFound {
		return err
	}

	// Exists already
	if err != nats.ErrStreamNotFound {
		return nil
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Retention: nats.LimitsPolicy,
	})

	if err != nil {
		return err
	}

	return nil
}

func (ex *NATSExchange) DeleteStream(streamName string) error {
	js, _ := ex.conn.JetStream()
	return js.DeleteStream(streamName)
}

func (ex *NATSExchange) Init(w Whisper) error {

	//TODO

	return nil
}

func (ex *NATSExchange) EmitSignal(subject string, payload []byte) error {
	return ex.conn.Publish(subject, payload)
}

func (ex *NATSExchange) EmitEvent(subject string, payload []byte) error {
	js, _ := ex.conn.JetStream()
	_, err := js.Publish(subject, payload)
	return err
}

func (ex *NATSExchange) SendBatchEvents(events []*Event) error {

	js, _ := ex.conn.JetStream()

	for _, ev := range events {

		m := &nats.Msg{
			Subject: ev.Subject,
			Data:    ev.Payload,
		}

		js.PublishMsg(m)
	}

	return ex.conn.Flush()
}

func (ex *NATSExchange) ChanSubscribe(subject string, channels chan Msg, durable string) (Subscription, error) {

	ch := make(chan *nats.Msg, 2048)

	js, _ := ex.conn.JetStream()

	// Durable
	var opts []nats.SubOpt
	if len(durable) > 0 {
		opts = append(opts, nats.Durable(durable))
	}

	sub, err := js.ChanSubscribe(subject, ch, opts...)

	sub.SetPendingLimits(-1, -1)
	ex.conn.Flush()

	s := &NATSSubscription{
		sub: sub,
		ch:  ch,
	}

	go func() {
		for m := range ch {
			channels <- &NATSMsg{
				m: m,
			}
		}
	}()

	return s, err
}

func (ex *NATSExchange) ChanWatch(subject string, channels chan Msg) (Subscription, error) {

	ch := make(chan *nats.Msg, 2048)

	sub, err := ex.conn.ChanSubscribe(subject, ch)

	sub.SetPendingLimits(-1, -1)
	ex.conn.Flush()

	s := &NATSSubscription{
		sub: sub,
		ch:  ch,
	}

	go func() {
		for m := range ch {
			channels <- &NATSMsg{
				m: m,
			}
		}
	}()

	return s, err
}
