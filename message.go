package whisper

import "encoding/json"

type Message struct {
	ID      string      `json:"id" binding:"required"`
	Type    string      `json:"type" binding:"required"` // normal, event
	Meta    Meta        `json:"meta" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

type Meta struct {
	Sender      *Member  `json:"from" binding:"required"`
	Group       string   `json:"group" binding:"required"`
	ContentType string   `json:"content_type"` // plain, image
	CreatedAt   int64    `json:"created_at" binding:"required"`
	Reference   *Message `json:"ref"`
}

type Member struct {
	ID          string `json:"id" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
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
