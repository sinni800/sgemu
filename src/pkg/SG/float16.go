package SG

import "math"

type FloatType byte

const (
	FloatViewRange = FloatType(0)
	FloatCD        = FloatType(1)
)

func Float16FromBits(n uint16) float32 {
	return float32(math.Floor(float64(float32(n)/2.0))) / 10.0
}

func Float16Bits(n float32) uint16 {
	return uint16(n * 20.0)
}

func Float16FromBits2(n uint16) float32 {
	return float32(math.Floor(float64(float32(n)/10.0))) / 10.0
}

func Float16Bits2(n float32) uint16 {
	return uint16(n * 100.0)
}
