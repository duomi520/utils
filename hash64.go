package utils

//代码抄自 https://github.com/zeebo/wyhash
import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

const (
	_wyp0 = 0xa0761d6478bd642f
	_wyp1 = 0xe7037ed1a0b428db
	_wyp2 = 0x8ebc6af09c88c6e3
	_wyp3 = 0x589965cc75374cc3
	_wyp4 = 0x1d8e4e27c47d124f
)

func _wymum(A, B uint64) uint64 {
	hi, lo := bits.Mul64(A, B)
	return hi ^ lo
}

func _wyr8(p unsafe.Pointer) uint64 {
	return binary.LittleEndian.Uint64((*[8]byte)(p)[:])
}

func _wyr4(p unsafe.Pointer) uint64 {
	return uint64(binary.LittleEndian.Uint32((*[4]byte)(p)[:]))
}

func _wyr3(p unsafe.Pointer, k uintptr) uint64 {
	b0 := uint64(*(*byte)(p))
	b1 := uint64(*(*byte)(offset(p, k>>1)))
	b2 := uint64(*(*byte)(offset(p, k-1)))
	return b0<<16 | b1<<8 | b2
}

func _wyr9(p unsafe.Pointer) uint64 {
	b := (*[8]byte)(p)
	return uint64(uint32(b[0])|uint32(b[1])<<8|uint32(b[2])<<16|uint32(b[3])<<24)<<32 |
		uint64(uint32(b[4])|uint32(b[5])<<8|uint32(b[6])<<16|uint32(b[7])<<24)
}

func hash(data string, seed uint64) uint64 {
	p, len := *(*unsafe.Pointer)(unsafe.Pointer(&data)), uintptr(len(data))
	see1, off := seed, len

	switch {
	case len <= 0x03:
		return _wymum(_wymum(_wyr3(p, len)^seed^_wyp0, seed^_wyp1)^seed, uint64(len)^_wyp4)

	case len <= 0x08:
		return _wymum(_wymum(_wyr4(offset(p, 0x00))^seed^_wyp0, _wyr4(offset(p, len-0x04))^seed^_wyp1)^seed, uint64(len)^_wyp4)

	case len <= 0x10:
		return _wymum(_wymum(_wyr9(offset(p, 0x00))^seed^_wyp0, _wyr9(offset(p, len-0x08))^seed^_wyp1)^seed, uint64(len)^_wyp4)

	case len <= 0x18:
		return _wymum(_wymum(_wyr9(offset(p, 0x00))^seed^_wyp0, _wyr9(offset(p, 0x08))^seed^_wyp1)^_wymum(_wyr9(offset(p, len-0x08))^seed^_wyp2, seed^_wyp3), uint64(len)^_wyp4)

	case len <= 0x20:
		return _wymum(_wymum(_wyr9(offset(p, 0x00))^seed^_wyp0, _wyr9(offset(p, 0x08))^seed^_wyp1)^_wymum(_wyr9(offset(p, 0x10))^seed^_wyp2, _wyr9(offset(p, len-0x08))^seed^_wyp3), uint64(len)^_wyp4)

	case len <= 0x100:
		seed = _wymum(_wyr8(offset(p, 0x00))^seed^_wyp0, _wyr8(offset(p, 0x08))^seed^_wyp1)
		see1 = _wymum(_wyr8(offset(p, 0x10))^see1^_wyp2, _wyr8(offset(p, 0x18))^see1^_wyp3)
		if len > 0x40 {
			seed = _wymum(_wyr8(offset(p, 0x20))^seed^_wyp0, _wyr8(offset(p, 0x28))^seed^_wyp1)
			see1 = _wymum(_wyr8(offset(p, 0x30))^see1^_wyp2, _wyr8(offset(p, 0x38))^see1^_wyp3)
			if len > 0x60 {
				seed = _wymum(_wyr8(offset(p, 0x40))^seed^_wyp0, _wyr8(offset(p, 0x48))^seed^_wyp1)
				see1 = _wymum(_wyr8(offset(p, 0x50))^see1^_wyp2, _wyr8(offset(p, 0x58))^see1^_wyp3)
				if len > 0x80 {
					seed = _wymum(_wyr8(offset(p, 0x60))^seed^_wyp0, _wyr8(offset(p, 0x68))^seed^_wyp1)
					see1 = _wymum(_wyr8(offset(p, 0x70))^see1^_wyp2, _wyr8(offset(p, 0x78))^see1^_wyp3)
					if len > 0xa0 {
						seed = _wymum(_wyr8(offset(p, 0x80))^seed^_wyp0, _wyr8(offset(p, 0x88))^seed^_wyp1)
						see1 = _wymum(_wyr8(offset(p, 0x90))^see1^_wyp2, _wyr8(offset(p, 0x98))^see1^_wyp3)
						if len > 0xc0 {
							seed = _wymum(_wyr8(offset(p, 0xa0))^seed^_wyp0, _wyr8(offset(p, 0xa8))^seed^_wyp1)
							see1 = _wymum(_wyr8(offset(p, 0xb0))^see1^_wyp2, _wyr8(offset(p, 0xb8))^see1^_wyp3)
							if len > 0xe0 {
								seed = _wymum(_wyr8(offset(p, 0xc0))^seed^_wyp0, _wyr8(offset(p, 0xc8))^seed^_wyp1)
								see1 = _wymum(_wyr8(offset(p, 0xd0))^see1^_wyp2, _wyr8(offset(p, 0xd8))^see1^_wyp3)
							}
						}
					}
				}
			}
		}

		off = (off-1)%0x20 + 1
		p = offset(p, len-off)

	default:
		for ; off > 0x100; off, p = off-0x100, offset(p, 0x100) {
			seed = _wymum(_wyr8(offset(p, 0x00))^seed^_wyp0, _wyr8(offset(p, 0x08))^seed^_wyp1) ^ _wymum(_wyr8(offset(p, 0x10))^seed^_wyp2, _wyr8(offset(p, 0x18))^seed^_wyp3)
			see1 = _wymum(_wyr8(offset(p, 0x20))^see1^_wyp1, _wyr8(offset(p, 0x28))^see1^_wyp2) ^ _wymum(_wyr8(offset(p, 0x30))^see1^_wyp3, _wyr8(offset(p, 0x38))^see1^_wyp0)
			seed = _wymum(_wyr8(offset(p, 0x40))^seed^_wyp0, _wyr8(offset(p, 0x48))^seed^_wyp1) ^ _wymum(_wyr8(offset(p, 0x50))^seed^_wyp2, _wyr8(offset(p, 0x58))^seed^_wyp3)
			see1 = _wymum(_wyr8(offset(p, 0x60))^see1^_wyp1, _wyr8(offset(p, 0x68))^see1^_wyp2) ^ _wymum(_wyr8(offset(p, 0x70))^see1^_wyp3, _wyr8(offset(p, 0x78))^see1^_wyp0)
			seed = _wymum(_wyr8(offset(p, 0x80))^seed^_wyp0, _wyr8(offset(p, 0x88))^seed^_wyp1) ^ _wymum(_wyr8(offset(p, 0x90))^seed^_wyp2, _wyr8(offset(p, 0x98))^seed^_wyp3)
			see1 = _wymum(_wyr8(offset(p, 0xa0))^see1^_wyp1, _wyr8(offset(p, 0xa8))^see1^_wyp2) ^ _wymum(_wyr8(offset(p, 0xb0))^see1^_wyp3, _wyr8(offset(p, 0xb8))^see1^_wyp0)
			seed = _wymum(_wyr8(offset(p, 0xc0))^seed^_wyp0, _wyr8(offset(p, 0xc8))^seed^_wyp1) ^ _wymum(_wyr8(offset(p, 0xd0))^seed^_wyp2, _wyr8(offset(p, 0xd8))^seed^_wyp3)
			see1 = _wymum(_wyr8(offset(p, 0xe0))^see1^_wyp1, _wyr8(offset(p, 0xe8))^see1^_wyp2) ^ _wymum(_wyr8(offset(p, 0xf0))^see1^_wyp3, _wyr8(offset(p, 0xf8))^see1^_wyp0)
		}
		for ; off > 0x20; off, p = off-0x20, offset(p, 0x20) {
			seed = _wymum(_wyr8(offset(p, 0x00))^seed^_wyp0, _wyr8(offset(p, 0x08))^seed^_wyp1)
			see1 = _wymum(_wyr8(offset(p, 0x10))^see1^_wyp2, _wyr8(offset(p, 0x18))^see1^_wyp3)
		}
	}

	switch {
	case off > 0x18:
		seed = _wymum(_wyr9(offset(p, 0x00))^seed^_wyp0, _wyr9(offset(p, 0x08))^seed^_wyp1)
		see1 = _wymum(_wyr9(offset(p, 0x10))^see1^_wyp2, _wyr9(offset(p, off-0x08))^see1^_wyp3)

	case off > 0x10:
		seed = _wymum(_wyr9(offset(p, 0x00))^seed^_wyp0, _wyr9(offset(p, 0x08))^seed^_wyp1)
		see1 = _wymum(_wyr9(offset(p, off-0x08))^see1^_wyp2, see1^_wyp3)

	case off > 0x08:
		seed = _wymum(_wyr9(offset(p, 0x00))^seed^_wyp0, _wyr9(offset(p, off-0x08))^seed^_wyp1)

	case off > 0x03:
		seed = _wymum(_wyr4(offset(p, 0x00))^seed^_wyp0, _wyr4(offset(p, off-0x04))^seed^_wyp1)

	default:
		seed = _wymum(_wyr3(p, off)^seed^_wyp0, seed^_wyp1)
	}

	return _wymum(seed^see1, uint64(len)^_wyp4)
}

// Hash64WY returns a 64bit digest of the data with different ones for every seed.
func Hash64WY[T string | []byte](data T, seed uint64) uint64 {
	if len(data) == 0 {
		return seed
	}
	return hash(*(*string)(unsafe.Pointer(&data)), seed)
}

var prime64 uint64 = 1099511628211

// FNV-1a算法
func Hash64FNV1A[T string | []byte](data T) uint64 {
	var result uint64 = 14695981039346656037
	for i := range len(data) {
		result ^= uint64(data[i])
		result *= prime64
	}
	return result
}

// https://github.com/wangyi-fudan/wyhash
// https://github.com/dgryski/go-wyhash/
