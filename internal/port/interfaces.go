package port

import (
	"context"
	"image"
)

type BookRepository interface {
	GetPage(ctx context.Context, bookID string, pageNum int) (Page, error)
	CurrentProgress(ctx context.Context, bookID string) (Progress, error)
	SaveProgress(ctx context.Context, bookID string, p Progress) error
}

type Renderer interface {
	Render(ctx context.Context, page Page, target Capabilities) (image.Image, error)
}

type Display interface {
	Capabilities() Capabilities
	Show(ctx context.Context, frame image.Image) error
	Clear(ctx context.Context) error
}

// InputSource emits events already tagged with which session they belong to.
type InputSource interface {
	Events(ctx context.Context) (<-chan Event, error)
}

type FrameCache interface {
	Get(key string) (image.Image, bool)
	Set(key string, frame image.Image)
}
