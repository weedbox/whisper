package whisper

type Member struct {
	ID          string `json:"id" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Avatar      string `json:"avatar"`
}

type GroupMember struct {
	Member   *Member `json:"member" validate:"required"`
	JoinedAt int64   `json:"joined_at"`
}
