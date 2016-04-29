package provisioner

import (
	"io"
	"sync"

	"github.com/Nanocloud/community/nanocloud/broadcaster"
)

type ProvFunc func(io.Writer)

type Provisioner struct {
	fn   ProvFunc
	cond *sync.Cond
	done bool

	b broadcaster.Broadcaster
}

func New(fn ProvFunc) *Provisioner {
	cond := sync.NewCond(&sync.Mutex{})

	return &Provisioner{
		fn:   fn,
		cond: cond,
	}
}

func (p *Provisioner) _run() {
	p.fn(&p.b)

	p.cond.L.Lock()
	p.done = true
	p.cond.Broadcast()
	p.cond.L.Unlock()
}
func (p *Provisioner) Run() {
	go p._run()
}

func (p *Provisioner) Wait() {
	p.cond.L.Lock()
	if !p.done {
		p.cond.Wait()
	}
	p.cond.L.Unlock()
}

func (p *Provisioner) AddOutput(w io.Writer) {
	p.b.Add(w)
}
