package Core

import "log"

type Func func()

type Runner struct {
	Funcs  chan Func
	Closer chan bool
}

func NewRunner() *Runner {
	return NewRunner2(10000)
}

func NewRunner2(size int) *Runner {
	f := new(Runner)
	f.Funcs = make(chan Func, size)
	f.Closer = make(chan bool)
	return f
}

func (r *Runner) Start() {
	go r.run()
}

func (r *Runner) Add(fnc Func) {
	r.Funcs <- fnc
}

func (r *Runner) Stop() {
	close(r.Funcs)
}

func (r *Runner) StopAndWait() {
	r.Stop()
	<-r.Closer
}

func (r *Runner) run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			go r.run()
		}
	}()
	for f := range r.Funcs {
		f()
	}
	r.Closer <- true
}
