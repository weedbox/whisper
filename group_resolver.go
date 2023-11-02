package whisper

import "errors"

var (
	ErrGroupNotFound    = errors.New("group_resolver: group not found")
	ErrOperationFailure = errors.New("group_resolver: operation failure")
)

type GroupResolver interface {
	Init(w Whisper) error
	GetMemberIDs(groupID string) ([]string, error)
	IsMutedMember(groupID string, memberID string) (bool, error)
}
