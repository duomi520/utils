package utils

import (
	"encoding/binary"
)

type ky struct {
	key   uint64
	value string
}

// MetaDict 非线程安全,key数量超过5个后，效率低于map
type MetaDict []ky

// GetAll
func (d MetaDict) GetAll() (s []string) {
	for _, k := range d {
		s = append(s, k.value)
	}
	return s
}

func (d MetaDict) findIndex(hash uint64) int {
	for i, k := range d {
		if k.key == hash {
			return i
		}
	}
	return -1
}

// Set 设置给定键的值，如果该键已存在，则更新值；如果不存在，则添加新的键值对。
func (d MetaDict) Set(key, value string) (new MetaDict) {
	hash := Hash64FNV1A(key)
	if idx := d.findIndex(hash); idx != -1 {
		d[idx].value = value
		return d
	}
	new = append(d, ky{hash, value})
	return new
}

// Get 根据给定的键返回相应的值。如果键存在，则返回对应的值和true；如果键不存在，则返回空字符串和false。
func (d MetaDict) Get(key string) (string, bool) {
	hash := Hash64FNV1A(key)
	if idx := d.findIndex(hash); idx != -1 {
		return d[idx].value, true
	}
	return "", false
}

// Del 根据给定的键删除相应的键值对。
func (d MetaDict) Del(key string) (new MetaDict) {
	hash := Hash64FNV1A(key)
	if idx := d.findIndex(hash); idx != -1 {
		if idx < (len(d) - 1) {
			copy(d[idx:], d[idx+1:])
		}
		return d[:len(d)-1]
	}
	return d
}

/*
+-------+-------+-------+-------+-------+-------+
| len(8)|        key (64)       |    value      |  ...
+-------+-------+-------+-------+-------+-------+
*/

// Encode 编码 将字典编码为字节切片。
func MetaDictEncode(d MetaDict) []byte {
	size := len(d)
	if size == 0 {
		return nil
	}
	n := 0
	for i := 0; i < size; i++ {
		n += 1 + 8 + len(d[i].value)
	}
	buf := make([]byte, n)
	idx := 0
	for i := 0; i < size; i++ {
		buf[idx] = byte(1 + 8 + len(d[i].value))
		idx++
		binary.LittleEndian.PutUint64(buf[idx:], d[i].key)
		idx += 8
		copy(buf[idx:], StringToBytes(d[i].value))
		idx += len(d[i].value)
	}
	return buf
}

// Decode 解码 将字节切片解码为字典。
func MetaDictDecode(data []byte) (d MetaDict) {
	idx := 0
	for len(data) > idx {
		n := int(data[idx])
		d = append(d, ky{binary.LittleEndian.Uint64(data[idx+1 : idx+1+8]), string(data[idx+1+8 : idx+n])})
		idx += n
	}
	return d
}
