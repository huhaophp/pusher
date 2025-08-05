package source

import (
	"context"
	"pusher/internal/types"
)

type Source interface {
	PullMessage(ctx context.Context, topic string, handler func(data *types.Data)) error
}
