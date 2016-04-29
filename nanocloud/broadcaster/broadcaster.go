package broadcaster

import (
	"io"
	"sync"
)

type Broadcaster struct {
	outs []io.Writer
	mut  sync.Mutex
}

func (b *Broadcaster) Add(w io.Writer) {
	b.mut.Lock()
	b.outs = append(b.outs, w)
	b.mut.Unlock()
}

func (b *Broadcaster) Write(buff []byte) (total int, err error) {
	b.mut.Lock()

	newOutLength := 0

	for i, out := range b.outs {
		if out != nil {
			written, e := out.Write(buff)
			total += written
			if e != nil {
				err = e
				b.outs[i] = nil
			} else {
				newOutLength++
			}
		}
	}

	if newOutLength != len(b.outs) {
		newOut := make([]io.Writer, newOutLength)
		i := 0
		for _, out := range b.outs {
			if out != nil {
				newOut[i] = out
				i++
			}
		}
		b.outs = newOut
	}

	b.mut.Unlock()

	return
}
