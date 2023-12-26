package commentary

import "os"

type Engine interface {
	Announce(speak string) (*os.File, error)
}
