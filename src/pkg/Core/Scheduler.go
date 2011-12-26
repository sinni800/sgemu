package Core

import (
	"container/list"
	"log"
	"time"
)

type Scheduler struct {
	Funcs  chan *Action
	List   *list.List
	Closer chan bool
}

type Action struct {
	fnc  Func
	Time time.Time
}

func NewScheduler() *Scheduler {
	return NewScheduler2(1000)
}

func NewScheduler2(size int) *Scheduler {
	f := new(Scheduler)
	f.Funcs = make(chan *Action, size)
	f.Closer = make(chan bool)
	f.List = list.New()
	return f
}

func (r *Scheduler) Test() {
	r.AddMS(func() { println("NOW") }, 0)
	r.AddMS(func() { println("5 ms") }, 5)
	r.AddSec(func() { println("5 sec") }, 5)
	r.AddSec(func() { println("10 sec") }, 10)
	r.AddSec(
		func() {
			println("2 sec and new add in 10 sec")
			r.AddSec(func() {
				println("12 sec and new add in 0 ms")
				r.AddMS(func() { println("now last") }, 1)
			}, 10)
		}, 2)
	r.AddTime(func() { println("time time 15 sec") }, time.Now().Add(time.Second*15))
	r.AddDur(func() { println("time dur 20 sec") }, time.Second*20)
}

func (r *Scheduler) Start() {
	go r.run()
}

func (r *Scheduler) AddTime(fnc Func, time time.Time) {
	r.Funcs <- &Action{fnc, time}
}

func (r *Scheduler) AddDur(fnc Func, dur time.Duration) {
	r.Funcs <- &Action{fnc, time.Now().Add(dur)}
}

func (r *Scheduler) AddMS(fnc Func, delayms time.Duration) {
	r.Funcs <- &Action{fnc, time.Now().Add(delayms * time.Millisecond)}
}

func (r *Scheduler) AddSec(fnc Func, delaysec time.Duration) {
	r.Funcs <- &Action{fnc, time.Now().Add(time.Second * delaysec)}
}

func (r *Scheduler) AddMin(fnc Func, delaysec time.Duration) {
	r.Funcs <- &Action{fnc, time.Now().Add(time.Minute * delaysec)}
}

func (r *Scheduler) Stop() {
	close(r.Funcs)
}

func (r *Scheduler) StopAndWait() {
	r.Stop()
	<-r.Closer
}

func (r *Scheduler) run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
			go r.run()
		}
	}()

	l := r.List

	dummyChan := make(<-chan time.Time)
	var ch = dummyChan

	var act *Action
	var elem *list.Element

	check := func() {
		for e := l.Front(); e != nil; e = e.Next() {
			a := e.Value.(*Action)
			if time.Now().After(a.Time) {
				a.fnc()
				l.Remove(e)
			} else {
				if act == nil {
					act = a
					elem = e
				} else if a.Time.Before(act.Time) {
					act = a
					elem = e
				}
			}
		}
		if act != nil {
			ch = time.After(act.Time.Sub(time.Now()))
		}
	}

	for {
		select {
		case act, ok := <-r.Funcs:
			if !ok {
				r.Closer <- true
				return
			}
			l.PushBack(act)
			check()
		case <-ch:
			ch = dummyChan
			act.fnc()
			act = nil
			l.Remove(elem)
			check()
		}
	}
}
