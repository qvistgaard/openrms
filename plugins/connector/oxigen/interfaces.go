package oxigen

import "io"

type ReadWriteCloserConnector interface {
	connect() (io.ReadWriteCloser, error)
}
