package whisper

import "encoding/json"

type Message struct {
	ID      string      `json:"id" validate:"required"`
	Type    string      `json:"type" validate:"required"` // normal, event
	Meta    Meta        `json:"meta" validate:"required"`
	Payload interface{} `json:"payload" validate:"required"`
}

type Meta struct {
	Sender      *Member  `json:"from" validate:"required"`
	Group       string   `json:"group" validate:"required"`
	ContentType string   `json:"content_type"` // plain, image
	CreatedAt   int64    `json:"created_at" validate:"required"`
	Reference   *Message `json:"ref"`
}

type Member struct {
	ID          string `json:"id" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
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
