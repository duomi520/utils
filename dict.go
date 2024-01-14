package utils

import (
	"encoding/binary"
)

// MetaDict 非线程安全,key数量超过5个后，效率低于map
type MetaDict struct {
	key   []uint64
	value []string
}

// Len 返回字典中键值对的数量
func (d *MetaDict) Len() int {
	return len(d.key)
}

// GetAll
func (d *MetaDict) GetAll() []string {
	return d.value
}

func (d *MetaDict) findIndex(hash uint64) int {
	for i, k := range d.key {
		if k == hash {
			return i
		}
	}
	return -1
}

// Set 设置给定键的值，如果该键已存在，则更新值；如果不存在，则添加新的键值对。
func (d *MetaDict) Set(key, value string) {
	hash := Hash64FNV1A(key)
	if idx := d.findIndex(hash); idx != -1 {
		d.value[idx] = value
		return
	}
	d.key = append(d.key, hash)
	d.value = append(d.value, value)
}

// Get 根据给定的键返回相应的值。如果键存在，则返回对应的值和true；如果键不存在，则返回空字符串和false。
func (d *MetaDict) Get(key string) (string, bool) {
	hash := Hash64FNV1A(key)
	if idx := d.findIndex(hash); idx != -1 {
		return d.value[idx], true
	}
	return "", false
}

// Del 根据给定的键删除相应的键值对。
func (d *MetaDict) Del(key string) {
	hash := Hash64FNV1A(key)
	if idx := d.findIndex(hash); idx != -1 {
		if idx < (d.Len() - 1) {
			copy(d.key[idx:], d.key[idx+1:])
			copy(d.value[idx:], d.value[idx+1:])
		}
		d.key = d.key[:d.Len()-1]
		d.value = d.value[:d.Len()-1]
		return
	}
}

/*
+-------+-------+-------+-------+-------+-------+
| len(8)|        key (64)       |    value      |  ...
+-------+-------+-------+-------+-------+-------+
*/

// Encode 编码 将字典编码为字节切片。
func (d *MetaDict) Encode() []byte {
	length := d.Len()
	if length == 0 {
		return nil
	}
	n := 0
	for i := 0; i < length; i++ {
		n += 1 + 8 + len(d.value[i])
	}
	buf := make([]byte, n)
	idx := 0
	for i := 0; i < length; i++ {
		buf[idx] = byte(1 + 8 + len(d.value[i]))
		idx++
		binary.LittleEndian.PutUint64(buf[idx:], d.key[i])
		idx += 8
		copy(buf[idx:], StringToBytes(d.value[i]))
		idx += len(d.value[i])
	}
	return buf
}

// Decode 解码 将字节切片解码为字典。
func (d *MetaDict) Decode(data []byte) {
	idx := 0
	for len(data) > idx {
		n := int(data[idx])
		d.key = append(d.key, binary.LittleEndian.Uint64(data[idx:idx+1+8]))
		d.value = append(d.value, string(data[idx+1+8:idx+n]))
		idx += n
	}
}
