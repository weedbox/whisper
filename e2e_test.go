package whisper

import (
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

var testNATSConn *nats.Conn
var testGM GroupManager

var members []Member
var groups []Group

func init() {
	nc, err := nats.Connect("0.0.0.0:32803")
	if err != nil {
		panic(err)
	}

	testNATSConn = nc

	// Initializing members
	for i := 0; i < 100; i++ {
		members = append(members, Member{
			ID:          fmt.Sprintf("member_%d", i),
			DisplayName: fmt.Sprintf("member_%d", i),
		})
	}

	// Initializing groups
	testGM = NewGroupManagerMemory()
	for i := 0; i < 100; i++ {

		var ms []string
		for _, m := range members {
			ms = append(ms, m.ID)
		}

		testGM.AddGroup(fmt.Sprintf("group_%d", i), ms)
	}
}

func Test_E2E(t *testing.T) {

	ex := NewNATSExchange(testNATSConn)

	w := NewWhisper(
		WithDomain("whisper"),
		WithBucketSize(32),
		WithExchange(ex),
		WithGroupManager(testGM),
	)

	err := w.Init()
	assert.Nil(t, err)

	groups := testGM.GetGroups()

	// Preparing a new message
	msg := Message{
		ID:   uuid.New().String(),
		Type: "normal",
		Meta: Meta{
			Sender:      &members[0],
			Group:       groups[0].ID(),
			ContentType: "plain",
		},
		Payload: "Hello",
	}

	// Initializing member agent to wait messages
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		ma, err := w.NewMemberAgent(members[0].ID)
		assert.Nil(t, err)

		ch := make(chan Msg, 2048)
		err = ma.ChanWatch(ch)
		assert.Nil(t, err)

		m := <-ch

		message, err := ParseMessage(m.Payload())
		assert.Nil(t, err)
		assert.Equal(t, msg.ID, message.ID)
		wg.Done()
	}()

	// Receive message
	mr, err := w.NewMessageReceiver()
	assert.Nil(t, err)
	defer mr.ClearStream()

	err = mr.Receive(msg.ToJSON())
	assert.Nil(t, err)

	// Dispatch to members
	md, err := w.NewMessageDispatcher()
	assert.Nil(t, err)
	defer md.ClearStream()

	var buckets []int32
	for i := int32(0); i < 32; i++ {
		buckets = append(buckets, i)
	}

	err = md.Subscribe(buckets)
	assert.Nil(t, err)

	wg.Wait()
}
