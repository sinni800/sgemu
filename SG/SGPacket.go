package SG

import (
	C "github.com/hjf288/sgemu/Core"
	"fmt"
	"io"
	"strings"
)

const (
	BufferSize = 1024
)

type SGPacket struct {
	C.Packet
	Encrypted bool
}

func NewPacket() (p *SGPacket) {
	p = new(SGPacket)
	p.Index = 0
	p.Buffer = make([]byte, BufferSize)
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

func (p *SGPacket) WriteColor(c Color) {
	p.WCheck(3)
	i := p.Index
	p.Buffer[i] = c.R
	p.Buffer[i+1] = c.G
	p.Buffer[i+2] = c.B
	p.Index = i + 3
}

func (p *SGPacket) ReadString() (pValue string) {
	return p.Packet.ReadString(int(p.ReadByte()))
}

func (p *SGPacket) WriteString(pValue string) {
	p.WriteByte(byte(len(pValue)))
	p.Packet.WriteString(pValue, len(pValue))
}

func (p *SGPacket) WriteHeader(opCode byte) {
	p.WCheck(5)
	p.WriteByte(0xAA)
	p.Index += 2
	p.WriteByte(opCode)
	p.Index++
}

func (p *SGPacket) ReadColor() Color {
	if !p.RCheck(3) {
		panic("Reading outside of the packet!")
	}
	i := p.Index
	p.Index = i + 3
	return NColor(p.Buffer[i], p.Buffer[i+1], p.Buffer[i+2])
}

func (p *SGPacket) ReadFloat(typ FloatType) float32 {
	i := p.ReadUInt16()
	f := float32(0)
	switch typ {
	case FloatViewRange:
		f = Float16FromBits(i)
		break
	case FloatCD:
		f = Float16FromBits2(i)
		break

	}
	return f
}

func (p *SGPacket) WriteFloat(f float32, typ FloatType) {
	i := uint16(0)
	switch typ {
	case FloatViewRange:
		i = Float16Bits2(f)
		break
	case FloatCD:
		i = Float16Bits2(f)
		break

	}
	p.WriteUInt16(i)
}

func (packet *SGPacket) ReadPacketFromStream(Reader io.Reader, callback func(*SGPacket)) (errno int) {
	bl, err := Reader.Read(packet.Buffer[packet.Index:])
	if err != nil {
		return -1
	}

	packet.Index += bl

	//enter when we recive enough bytes to start reading them
	for packet.Index > 2 {
		p := packet
		size := p.Index
		p.Index = 0

		//Check header byte
		if p.ReadByte() != 0xAA {
			//AA == SG Packet 
			panic("Wrong packet header")
			//client.Log().Printf("Wrong packet header")
			//client.Log().Printf("% #X", p.Buffer[:size])
			return -2
		}
		l := int(p.ReadUInt16())
		p.Index = size

		if len(packet.Buffer) < l {
			packet.Resize(l)
		}

		//enter when we recived enough packet data
		if size >= l+3 {
			temp := packet.Buffer[:l+3]
			op := packet.Buffer[3]

			//check if packet is encrypted
			if op > 13 || (op > 1 && op < 5) || (op > 6 && op < 13) {
				var sumCheck bool
				temp, sumCheck = DecryptPacket(temp)
				if !sumCheck {
					panic("Packet sum check failed!")
					return -3
				}
			} else {
				temp = temp[3:]
			}

			//handle packet
			callback(NewPacketRef(temp))
			packet.Index = 0

			//enter when we have more than one packet in buffer
			if size > l+3 {
				packet.Index = size - (l + 3)
				copy(packet.Buffer, packet.Buffer[l+3:size])
			} else {
				//enter when we done processing the buffer
				//keeping the user under 4048k use to save memory
				if cap(packet.Buffer) > 4048 {
					packet.Buffer = make([]byte, BufferSize)
					packet.Index = 0
				}
			}
		} else {
			//break if we didn't get all the packet bytes
			break
		}
	}
	return 0
}

func (p *SGPacket) String() string {
	pString := strings.Replace(string(p.Buffer), "\a", " ", -1)
	return fmt.Sprintf("Header(%d) len(%d) : % #X\n %s", p.Buffer[0], len(p.Buffer), p.Buffer, pString)
}
