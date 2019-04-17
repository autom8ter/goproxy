package util

import (
	"github.com/autom8ter/objectify"
	"io"
)

var Handle = objectify.Default()

// switchProtocolCopier exists so goroutines proxying data back and
// forth have nice names in stacks.
type SwitchProtocolCopier struct {
	User, Backend io.ReadWriter
}

func (c SwitchProtocolCopier) CopyFromBackend(errc chan<- error) {
	_, err := io.Copy(c.User, c.Backend)
	errc <- err
}

func (c SwitchProtocolCopier) CopyToBackend(errc chan<- error) {
	_, err := io.Copy(c.Backend, c.User)
	errc <- err
}
