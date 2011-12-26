package SG

var (
	Key1 []byte
	Key2 []byte
) 

func init() {
	sKey1 := "NexonInc.NexonInc.NexonInc.NexonInc."
	Key1 = make([]byte, len(sKey1))
	copy(Key1, sKey1)
	Key2 = make([]byte, 256*4)
   
	for i := 0; i < 256; i++ {
		Key2[i*4] = uint8(i)
		Key2[(i*4)+1] = uint8(i)
		Key2[(i*4)+2] = uint8(i)
		Key2[(i*4)+3] = uint8(i)
	} 
}
 

func DecryptPacket(bufferf []byte) ([]byte, bool) {
	buffer := bufferf[:len(bufferf)-2]

	if len(buffer) < 6 {
		return buffer, true
	} 

	cSum := byte(0)
	for i := 0; i < len(bufferf[3:])-1; i++ {
		cSum += bufferf[3+i]
	}

	if cSum != 0 {
		return buffer, false
	}

	remider := len(buffer[5:]) % 9
	it := (len(buffer[5:]) - remider) / 9

	keypos := int(buffer[4]) * 4
	Crypt(buffer[5:], Key2[keypos:keypos+4])
	for i := 0; i < it; i++ {
		if byte(i) != buffer[4] {
			keypos = (i % 256) * 4
			Crypt(buffer[5+(i*9):5+(i*9)+9], Key2[keypos:keypos+4])
		}
	}

	if remider > 0 && byte(it) != buffer[4] {
		keypos = (it % 256) * 4
		Crypt(buffer[5+(it*9):5+(it*9)+remider], Key2[keypos:keypos+4])
	}
	Crypt(buffer[5:], Key1[0:36])
	buffer[4] = buffer[3]

	cSum = byte(0)
	for i := 0; i < len(buffer[4:]); i++ {
		cSum += buffer[4+i]
	}
	cSum += bufferf[len(bufferf)-1]

	if cSum != 0 {
		return buffer[4:len(buffer)], false
	}

	return buffer[4:len(buffer)], true
}

func Crypt(bufferSliced []byte, KeySliced []byte) {
	for i := 0; i < len(bufferSliced); i++ {
		bufferSliced[i] ^= KeySliced[i%len(KeySliced)]
	}
}

func EncryptPacket(bufferf []byte, packetNum uint8) []byte {
	buffer := bufferf[:len(bufferf)-2]

	buffer[4] = packetNum

	if len(buffer) < 6 {
		return buffer
	}

	cSum := buffer[3]
	for i := 0; i < len(buffer[5:]); i++ {
		cSum += buffer[5+i]
	}
	cSum = 255 - cSum
	bufferf[len(bufferf)-1] = cSum + 1

	remider := len(buffer[5:]) % 9
	it := (len(buffer[5:]) - remider) / 9

	Crypt(buffer[5:], Key1[0:36])
	keypos := int(buffer[4]) * 4
	Crypt(buffer[5:], Key2[keypos:keypos+4])
	for i := 0; i < it; i++ {
		if byte(i) != packetNum {
			keypos = (i % 256) * 4
			Crypt(buffer[5+(i*9):5+(i*9)+9], Key2[keypos:keypos+4])
		}
	}

	if remider > 0 && byte(it) != packetNum {
		keypos = (it % 256) * 4
		Crypt(buffer[5+(it*9):5+(it*9)+remider], Key2[keypos:keypos+4])
	}

	cSum = byte(0)
	for i := 0; i < len(buffer[3:]); i++ {
		cSum += buffer[3+i]
	}
	cSum = 255 - cSum
	bufferf[len(bufferf)-2] = cSum + 1

	return buffer
}
