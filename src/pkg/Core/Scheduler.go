package Core

import (
	"time"
	"container/list"
)

type Scheduler struct {
	Funcs chan *Action
	List *list.List
}

type Action struct {
	fnc Func
	Time *time.Time
}

func NewScheduler() *Scheduler {
	return NewScheduler2(1000)
} 

func NewScheduler2(size int) *Scheduler {
	f := new(Scheduler)
	f.Funcs = make(chan *Action, size)
	f.List = list.New()
	return f
}


func (r *Scheduler) Start() {
	go r.run()
}

func (r *Scheduler) Add(fnc Func, time *time.Time) {
	r.Funcs <- &Action{fnc,time}
}

func (r *Scheduler) AddMS(fnc Func, delayms int64) {
	t := time.Now().Add(time.Duration(delayms*1000*1000));
	r.Funcs <- &Action{fnc,&t}
}

func (r *Scheduler) AddSec(fnc Func, delaysec int64) {
	t := time.Now().Add(time.Duration(delaysec*1000*1000*1000));
	r.Funcs <- &Action{fnc,&t}
}

func (r *Scheduler) Stop() {
	close(r.Funcs)
}

func (r *Scheduler) run() {
	l := r.List
	var t *time.Time = nil
	
	dummyChan := make(<-chan time.Time)
	var ch = dummyChan
	
	var fnc Func
	var elem *list.Element
	
	check := func() {
		for e := l.Front(); e != nil; e = e.Next() {
			a := e.Value.(*Action)
			if time.Now().After(*a.Time) {
				a.fnc()
				l.Remove(e)
			} else {
				if (t == nil) {
					t = a.Time
					fnc = a.fnc
					elem = e
				} else if (a.Time.Before(*t)) {
					t = a.Time
					fnc = a.fnc
					elem = e
				}
			}
		}
		if t != nil { 
			ch = time.After(t.Sub(time.Now())) 
		}	
	} 
	
	for { 
		select {
			case act := <-r.Funcs:
				l.PushBack(act)	
				check()
			case <-ch: 
				ch = dummyChan
				t = nil
				fnc()
				l.Remove(elem)
				check()
		}
	}
}
