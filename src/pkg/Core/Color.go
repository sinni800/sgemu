package Core

var (
	Red       = NColor(255, 0, 0)
	Green     = NColor(0, 255, 0)
	Blue      = NColor(0, 0, 255)
	Black     = NColor(0, 0, 0)
	White     = NColor(255, 255, 255)
	HelpColor = NColor(0x46, 0xFA, 0xC8)
) 

type Color struct {
	R, G, B byte
}

func NColor(r, g, b byte) *Color {
	return &Color{r, g, b}
}
