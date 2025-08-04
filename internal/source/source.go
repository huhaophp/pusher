package source

import (
	"context"
	"pusher/internal/model"
)

type Source interface {
	PullMessage(ctx context.Context, topic string, handler func(data *model.Data)) error
}
