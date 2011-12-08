package Core
/*
import "math"

type Packet struct {
	Buffer []byte
	Index  int
}

func NewPacket() (p *Packet) {
	p = new(Packet)
	p.Index = 0
	p.Buffer = make([]byte, 1024+7)
	return p
}

func NewPacket2(size int) (p *Packet) {
	p = new(Packet)
	p.Index = 0
	p.Buffer = make([]byte, size+7)
	return p
}

func NewPacket3(buffer []byte) (p *Packet) {
	p = new(Packet)
	p.Buffer = make([]byte, len(buffer))
	p.Index = 0
	copy(p.Buffer, buffer)
	return p
}

func NewPacketRef(buffer []byte) (p *Packet) {
	p = new(Packet)
	p.Buffer = buffer
	p.Index = 0
	return p
}

func (p *Packet) WCheck(size int) bool {
	if p.Index+size > len(p.Buffer) {
		return false
	}
	return true
}

func (p *Packet) RCheck(size int) bool {
	if p.Index+size > len(p.Buffer) {
		return false
	}
	return true
}

func (p *Packet) Resize(newSize int) {
	if len(p.Buffer) > newSize {
		p.Buffer = p.Buffer[:newSize]
		return
	}
	temp := make([]byte, newSize)
	copy(temp, p.Buffer)
	p.Buffer = temp
}

func (p *Packet) WSkip(times int) {
	if !p.WCheck(times) {
		return
	}
	p.Index += times
}

func (p *Packet) RSkip(times int) {
	if !p.RCheck(times) {
		return
	}
	p.Index += times
}

func (p *Packet) Write(bytes []byte) {
	if !p.WCheck(len(bytes)) {
		return
	}
	copy(p.Buffer[p.Index:], bytes)
	p.Index += len(bytes)
}

func (p *Packet) WriteByte(b byte) {
	if !p.WCheck(1) {
		return
	}
	p.Buffer[p.Index] = b
	p.Index++
}

func (p *Packet) WriteHeader(opCode byte) {
	if !p.WCheck(5) {
		return
	}
	p.WriteByte(0xAA)
	p.Index += 2
	p.WriteByte(opCode)
	p.Index++
}


func (p *Packet) WriteLen() {
	t := p.Index
	p.Index = 1
	p.WriteUInt16(uint16(t - 3))
	p.Buffer[1], p.Buffer[2] = p.Buffer[2], p.Buffer[1]
	p.Index = t
}


func (p *Packet) WriteInt16(value int16) {
	if !p.WCheck(2) {
		return
	}
	p.Buffer[p.Index] = (byte)(value)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 8)
	p.Index++
}

func (p *Packet) WriteUInt16(value uint16) {
	if !p.WCheck(2) {
		return
	}
	p.Buffer[p.Index] = (byte)(value)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 8)
	p.Index++
}

func (p *Packet) WriteInt32(value int32) {
	if !p.WCheck(4) {
		return
	}
	p.Buffer[p.Index] = (byte)(value)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 8)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 16)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 24)
	p.Index++
}

func (p *Packet) WriteUInt32(value uint32) {
	if !p.WCheck(4) {
		return
	}
	p.Buffer[p.Index] = (byte)(value)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 8)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 16)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 24)
	p.Index++
}

func (p *Packet) WriteInt64(value int64) {
	if !p.WCheck(8) {
		return
	}
	p.Buffer[p.Index] = (byte)(value)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 8)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 16)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 24)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 32)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 40)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 48)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 56)
	p.Index++
}

func (p *Packet) WriteUInt64(value uint64) {
	if !p.WCheck(8) {
		return
	}
	p.Buffer[p.Index] = (byte)(value)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 8)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 16)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 24)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 32)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 40)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 48)
	p.Index++
	p.Buffer[p.Index] = (byte)(value >> 56)
	p.Index++
}


func (p *Packet) WriteFloat32(pValue float32) {
	if !p.WCheck(4) {
		return
	}
	p.WriteUInt32(math.Float32bits(pValue))
}

func (p *Packet) WriteFloat64(pValue float64) {
	if !p.WCheck(8) {
		return
	}
	p.WriteUInt64(math.Float64bits(pValue))
}

func (p *Packet) WriteString(pValue string) {
	p.WriteByte(byte(len(pValue)))
	p.oWriteString(pValue, len(pValue))
}

func (p *Packet) oWriteString(pValue string, size int) {
	if !p.WCheck(size) {
		return
	}
	copy(p.Buffer[p.Index:], []byte(pValue))
	p.Index += size
}

func (p *Packet) WriteRawString(pValue string) {
	p.oWriteString(pValue, len(pValue))
}

func (p *Packet) ReadBytes(size int) (pValue []byte) {
	if !p.RCheck(size) {
		return
	}
	pValue = p.Buffer[p.Index : p.Index+size]
	p.Index += size
	return pValue
}

func (p *Packet) ReadByte() (pValue byte) {
	if !p.RCheck(1) {
		return
	}
	pValue = byte(p.Buffer[p.Index])
	p.Index++
	return pValue
}

func (p *Packet) ReadInt16() (pValue int16) {
	if !p.RCheck(2) {
		return
	}
	pValue = int16(p.Buffer[p.Index])
	p.Index++
	pValue |= (int16(p.Buffer[p.Index]) << 8)
	p.Index++
	return pValue
}

func (p *Packet) ReadUInt16() (pValue uint16) {
	if !p.RCheck(2) {
		return
	}
	pValue = uint16(p.Buffer[p.Index])
	p.Index++
	pValue |= (uint16(p.Buffer[p.Index]) << 8)
	p.Index++
	return pValue
}

func (p *Packet) ReadInt32() (pValue int32) {
	if !p.RCheck(4) {
		return
	}
	pValue = int32(p.Buffer[p.Index])
	p.Index++
	pValue |= (int32(p.Buffer[p.Index]) << 8)
	p.Index++
	pValue |= (int32(p.Buffer[p.Index]) << 16)
	p.Index++
	pValue |= (int32(p.Buffer[p.Index]) << 24)
	p.Index++
	return pValue
}

func (p *Packet) ReadUInt32() (pValue uint32) {
	if !p.RCheck(4) {
		return
	}
	pValue = uint32(p.Buffer[p.Index])
	p.Index++
	pValue |= (uint32(p.Buffer[p.Index]) << 8)
	p.Index++
	pValue |= (uint32(p.Buffer[p.Index]) << 16)
	p.Index++
	pValue |= (uint32(p.Buffer[p.Index]) << 24)
	p.Index++
	return pValue
}

func (p *Packet) ReadInt64() (pValue int64) {
	if !p.RCheck(4) {
		return
	}
	pValue = int64(p.Buffer[p.Index])
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 8)
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 16)
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 24)
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 32)
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 40)
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 48)
	p.Index++
	pValue |= (int64(p.Buffer[p.Index]) << 56)
	p.Index++
	return pValue
}

func (p *Packet) ReadUInt64() (pValue uint64) {
	if !p.RCheck(8) {
		return
	}
	pValue = uint64(p.Buffer[p.Index])
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 8)
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 16)
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 24)
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 32)
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 40)
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 48)
	p.Index++
	pValue |= (uint64(p.Buffer[p.Index]) << 56)
	p.Index++
	return pValue
}

func (p *Packet) ReadFloat32() (pValue float32) {
	if !p.RCheck(4) {
		return
	}
	pValue = math.Float32frombits(p.ReadUInt32())
	return pValue
}

func (p *Packet) ReadFloat64() (pValue float64) {
	if !p.RCheck(8) {
		return
	}
	pValue = math.Float64frombits(p.ReadUInt64())
	return pValue
}

func (p *Packet) ReadString() (pValue string) {
	return p.oReadString(int(p.ReadByte()))
}


func (p *Packet) oReadString(size int) (pValue string) {
	if !p.RCheck(size) {
		return
	}
	pValue = string(p.Buffer[p.Index : p.Index+size])
	p.Index += size
	return pValue
}
*/
