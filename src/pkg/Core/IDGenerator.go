package Core

type IDGen struct {
	cGen       chan uint32
	lastNumber uint32
	limited    bool
}

func NewIDG() *IDGen {
	return NewIDG3(3000, true)
}

func NewIDG2(size int) *IDGen {
	return NewIDG3(size, false)
}

func NewIDG3(size int, limited bool) *IDGen {
	g := &IDGen{ make(chan uint32, size), 1, limited}
	go g.Gen()
	return g
}
  
func (g *IDGen) Gen() {
	for i := 0; i < cap(g.cGen); i++ {
		g.cGen <- g.lastNumber
		g.lastNumber++
	}
}

func (g *IDGen) Next() (id uint32, full bool) {
	select {
	case id = <-g.cGen:
	default:
		return 0, false
	}
	return id, true
}

func (g *IDGen) Return(id uint32) {
	g.cGen <- id
}
