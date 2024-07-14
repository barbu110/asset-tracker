package rendered_label_storage

import (
	"io"
	"time"
)

type RenderedLabelStorage interface {
	Create(key string, data io.Reader) error
	GetDownloadUrl(key string, validity time.Duration) (string, error)
}
