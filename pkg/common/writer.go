package common

import "io"

type WriterWrapper interface {
	Writer() io.Writer
}
