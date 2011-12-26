package SG

import (
	C "Core"
	"fmt"
)

type Packet interface {
	BasePacket() *C.Packet
} 

type SGPacket struct {
	 C.Packet 
}
 
func NewPacket() (p *SGPacket) {
	p = new(SGPacket)
	p.Index = 0
	p.Buffer = make([]byte, 1024)
	return p
}

func NewPacket2(size int) (p *SGPacket) {
	p = new(SGPacket)
	p.Index = 0
	p.Buffer = make([]byte, size)
	return p
} 
 
func NewPacket3(buffer []byte) (p *SGPacket) {
	p = new(SGPacket)
	p.Buffer = make([]byte, len(buffer))
	p.Index = 0
	copy(p.Buffer, buffer)
	return p
}

func NewPacketRef(buffer []byte) (p *SGPacket) {
	p = new(SGPacket)
	p.Buffer = buffer
	p.Index = 0
	return p
}
 
func (p *SGPacket) WriteColor(c *Color) {
	p.WCheck(3)
	i := p.Index
	p.Buffer[i] = c.R
	p.Buffer[i+1] = c.G
	p.Buffer[i+2] = c.B 
	p.Index = i+3
} 

func (p *SGPacket) ReadString() (pValue string) {
	return p.Packet.ReadString(int(p.ReadByte()))
}

func (p *SGPacket) WriteString(pValue string) {
	p.WriteByte(byte(len(pValue)))
	p.Packet.WriteString(pValue, len(pValue))
}

func (p *SGPacket) WriteHeader(opCode byte) {
	if !p.WCheck(5) {
		return
	}
	p.WriteByte(0xAA)
	p.Index += 2
	p.WriteByte(opCode)
	p.Index++
} 

func (p *SGPacket) ReadColor() *Color{
	if !p.RCheck(3) {
		panic("Reading outside of the packet!")
	}
	i := p.Index
	p.Index = i+3
	return NColor(p.Buffer[i],p.Buffer[i+1],p.Buffer[i+2])
} 

func (p *SGPacket) String() (string) { 
	return fmt.Sprintf("Header(%d) len(%d) : % #X\n %s", p.Buffer[0], len(p.Buffer), p.Buffer, p.Buffer)
}