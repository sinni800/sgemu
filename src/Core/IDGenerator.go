package Core

type IDGen struct {
	cGen       chan uint32
	lastNumber uint32
}

func NewIDG() *IDGen {
	return NewIDG3(3000, 1)
}

func NewIDG2(size int) *IDGen {
	return NewIDG3(size, 1)
}

func NewIDG3(size int, startNum uint32) *IDGen {
	g := &IDGen{ make(chan uint32, size), startNum}
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
