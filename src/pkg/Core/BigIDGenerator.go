package Core

type BIDGen struct {
	cGen       chan uint32
	cHelper    chan uint32
	cSignal    chan bool
	tempHeap   []uint32
	lastNumber uint32
}

func NewBIDG() *BIDGen {
	return NewBIDG2(1000)
} 

func NewBIDG2(size int) *BIDGen {
	g := &BIDGen{make(chan uint32, size),make(chan uint32, 100),make(chan bool, 10),make([]uint32, 0),0}
	go g.Gen()
	return g
}

func (g *BIDGen) Gen() {
	g.cSignal <- true
	for {
	Signal:
		switch <-g.cSignal {
		case true:
			for {
				if len(g.tempHeap) > 0 {
					id := g.tempHeap[len(g.tempHeap)-1]
					select {
					case g.cGen <- id:
						g.tempHeap = g.tempHeap[:len(g.tempHeap)-1]
					default:
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
					g.tempHeap = append(g.tempHeap, id)
				default:
					goto Signal
				}
			}
		}
	}
}


func (g *BIDGen) Next() (id uint32) {
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

func (g *BIDGen) Return(id uint32) {
	select {
	case g.cGen <- id:
	default:
		select {
			case g.cHelper <- id:
				g.cSignal <- false
			default:
				go func() {  g.cHelper <- id;  g.cSignal <- false   }()
			}
	}
}
