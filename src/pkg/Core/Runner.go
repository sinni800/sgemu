package Core

type Func func()

type Runner struct {
	Funcs chan Func
}

func NewRunner() *Runner {
	return NewRunner2(10000)
}

func NewRunner2(size int) *Runner {
	f := new(Runner)
	f.Funcs = make(chan Func, size)
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

func (r *Runner) run() {
	for f := range r.Funcs {
		f()
	}
}
