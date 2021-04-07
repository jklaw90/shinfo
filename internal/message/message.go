package message

import (
	"context"

	"github.com/jklaw90/shinfo/pkg/model"
)

type Service interface {
	List(ctx context.Context, params ListParams) (model.MessageList, error)
	AddMessage(ctx context.Context, msg model.MessageCreate) (model.Message, error)
	// Get(ctx context.Context, roomID, messageID gocql.UUID) (model.Message, error)
	// Delete(ctx context.Context, roomID, messageID gocql.UUID) error
}
