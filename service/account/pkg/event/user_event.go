package event

import (
	"github.com/abdelrahman146/zard/service/account/pkg/usecase"
	"github.com/abdelrahman146/zard/shared"
	"github.com/abdelrahman146/zard/shared/pubsub/messages"
)

func SubscribeToUserEvents(usecases *usecase.AccountUseCases, toolkit shared.Toolkit) error {
	ue := &userEvent{
		toolkit:  toolkit,
		usecases: usecases,
	}
	var err error
	if _, err = toolkit.PubSub.Subscribe(&messages.UserCreatedMessage{}, ue.UserCreated); err != nil {
		return err
	}
	return nil
}

type userEvent struct {
	toolkit  shared.Toolkit
	usecases *usecase.AccountUseCases
}

func (e *userEvent) UserCreated(received []byte) error {
	return nil
}
