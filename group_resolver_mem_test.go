package whisper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GroupResolverMemory_GetMemverIDs(t *testing.T) {

	gr := NewGroupResolverMemory()

	gr.addGroup("test", []string{
		"fred",
		"jhe",
		"chuck",
		"leon",
	})

	ids, err := gr.GetMemberIDs("test")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(ids))
}

func Test_GroupResolverMemory_IsMutedMember(t *testing.T) {
	gr := NewGroupResolverMemory()

	gr.addGroup("test", []string{
		"fred",
	})

	muted, err := gr.IsMutedMember("test", "fred")
	assert.Nil(t, err)
	assert.False(t, muted)

	gr.addMutedMembers("test", []string{
		"fred",
	})

	muted, err = gr.IsMutedMember("test", "fred")
	assert.Nil(t, err)
	assert.True(t, muted)
}
