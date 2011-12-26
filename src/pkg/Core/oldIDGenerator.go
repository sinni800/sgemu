package Core

/*
import (
	H "container/vector"
)

type IDGen struct {
	cGen       chan uint32
	cHelper    chan uint32
	tempHeap   *H.Vector
	cSignal    chan bool
	lastNumber uint32
}

func NewIDG() *IDGen {
	return NewIDG2(1000)
} 

func NewIDG2(size int) *IDGen {
	g := new(IDGen)
	g.cGen = make(chan uint32, size)
	g.cHelper = make(chan uint32, size)
	g.cSignal = make(chan bool, 10)
	g.tempHeap = new(H.Vector)
	g.lastNumber = 0
	go g.Gen()
	return g
}

func (g *IDGen) Gen() {
	g.cSignal <- true
	for {
	Signal:
		switch <-g.cSignal {
		case true:
			for {
				if g.tempHeap.Len() > 0 {
					id := g.tempHeap.Pop().(uint32)
					select {
					case g.cGen <- id:
					default:
						g.tempHeap.Push(id)
						goto Signal
					}
				}

				select {
				case g.cGen <- g.lastNumber:
					g.lastNumber++
				default:
					goto Signal
				}
			}

		case false:
			id := uint32(0)
			for {
				select {
				case id = <-g.cHelper:
					g.tempHeap.Push(id)
				default:
					goto Signal
				}
			}
		}
	}
}


func (g *IDGen) Next() (id uint32) {
	select {
	case id = <-g.cGen:
	default:

		// old way, faster but wont last forever. can be used though.
		//id = g.lastNumber	
		//g.lastNumber++
		//g.cSignal <- true


		//little slower but most stable way
		g.cSignal <- true
		id = <-g.cGen
	}
	return id
}

func (g *IDGen) Return(id uint32) {
	select {
	case g.cGen <- id:
	default:
		g.cHelper <- id
		g.cSignal <- false
	}
}

*/
