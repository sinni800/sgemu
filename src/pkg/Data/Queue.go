package Data

import . "container/list"
import "time"

var LoginQueue *Queue = NewQueue(false)

type Item interface{}

type Queue struct {
	In      chan *InStruct
	Request chan chan Item
	List    *List
	IPCheck bool
}

type InStruct struct {
	IP   string
	ID   string
	Time time.Time
}

func NewQueue(ipcheck bool) *Queue {
	q := new(Queue)
	q.In = make(chan *InStruct, 50)
	q.Request = make(chan chan Item, 50)
	q.List = New()
	q.IPCheck = ipcheck
	go q.run()
	return q
}  

func (q *Queue) Add(ip, id string) {
	q.In <- &InStruct{ip, id, time.Time{}}
}
  
func (q *Queue) Check(ip string) (string, bool) {
	chn := make(chan Item)
	q.Request <- chn
	chn <- ip
	id := <-chn

	if id == nil {
		return "", false
	}
	return id.(string), true
}

func (q *Queue) run() {
	timer := time.After(10 * 1e9)
	for {
	SELECT:
		select {

		case i := <-q.In:
			i.Time = time.Now()

			if q.IPCheck {
				if q.List.Len() > 0 {
					for e := q.List.Front(); e != nil; e = e.Next() {
						if v, ok := e.Value.(*InStruct); ok {
							if v.IP == i.IP {
								goto SELECT
							}
						}
					}
				}
			}

			q.List.PushFront(i)

		case r := <-q.Request:
			ip := (<-r).(string)
			for e := q.List.Front(); e != nil; e = e.Next() {
				if i, ok := e.Value.(*InStruct); ok {
					if i.IP == ip {
						q.List.Remove(e)
						r <- i.ID
						goto SELECT
					}
				}
			}
			r <- nil

		case <-timer:
			for e := q.List.Front(); e != nil; e = e.Next() {
				if i, ok := e.Value.(*InStruct); ok {
					if (time.Now().Unix() - i.Time.Unix()) >= 10 {
						q.List.Remove(e) 
					}
				}
			}
			timer = time.After(10 * 1e9)
		}
	}
}
