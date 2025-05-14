package platform

import "context"

type AudioPlayer interface {
	PlayLoop(ctx context.Context, fileName string)
}
