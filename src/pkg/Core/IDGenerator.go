package Core

type IDGen struct {
	cGen       chan uint32
	lastNumber uint32
}

func NewIDG() *IDGen {
	return NewIDG2(3000)
}

func NewIDG2(size int) *IDGen {
	g := new(IDGen)
	g.cGen = make(chan uint32, size)
	g.lastNumber = 1
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
