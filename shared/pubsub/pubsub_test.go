package pubsub

import (
	"encoding/json"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestPubsub(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := NewMockPubSub(ctrl)
	activity := &messages.NewActivity{
		AppID: "test",
	}
	m.EXPECT().Publish(activity).Return(nil)
	m.EXPECT().Subscribe(activity, gomock.Any()).Return(nil, nil)
	err := m.Publish(activity)
	assert.NoError(t, err)
	_, err = m.Subscribe(activity, func(received []byte) error {
		a := &messages.NewActivity{}
		err := json.Unmarshal(received, a)
		return err
	})
	assert.NoError(t, err)

}
