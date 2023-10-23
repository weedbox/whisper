package main

import "encoding/json"

type Message struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"` // normal, event
	Meta    Meta        `json:"meta"`
	Payload interface{} `json:"payload"`
}

type Meta struct {
	Sender      *Member  `json:"from"`
	Group       string   `json:"group"`
	ContentType string   `json:"content_type"` // plain, image
	CreatedAt   int64    `json:"created_at"`
	Reference   *Message `json:"ref"`
}

type Member struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

func ParseMessage(data []byte) (*Message, error) {

	msg := &Message{}

	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, ErrInvalidMessageFormat
	}

	return msg, nil
}

func (m *Message) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}
