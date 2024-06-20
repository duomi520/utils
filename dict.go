package utils

import (
	"encoding/binary"
	"slices"
)

type iMetaDict interface {
	string | any
}

// MetaDict 非线程安全,key数量超过5个后，效率低于map
type MetaDict[T iMetaDict] struct {
	key   []uint64
	value []T
}

// Len 长度
func (m MetaDict[iMetaDict]) Len() int {
	return len(m.key)
}

// GetAll
func (m MetaDict[iMetaDict]) GetAll() (s []iMetaDict) {
	return m.value
}

// Set 设置给定键的值，如果该键已存在，则更新值；如果不存在，则添加新的键值对。
func (m MetaDict[iMetaDict]) Set(key string, value iMetaDict) MetaDict[iMetaDict] {
	hash := Hash64FNV1A(key)
	if idx := slices.Index(m.key, hash); idx > -1 {
		m.value[idx] = value
		return m
	}
	m.key = append(m.key, hash)
	m.value = append(m.value, value)
	return m
}

// Get 根据给定的键返回相应的值。如果键存在，则返回对应的值和true；如果键不存在，则返回空字符串和false。
func (m MetaDict[iMetaDict]) Get(key string) (v iMetaDict, ok bool) {
	hash := Hash64FNV1A(key)
	if idx := slices.Index(m.key, hash); idx > -1 {
		v = m.value[idx]
		ok = true
		return
	}
	ok = false
	return
}

// Del 根据给定的键删除相应的键值对。
func (m MetaDict[iMetaDict]) Del(key string) MetaDict[iMetaDict] {
	hash := Hash64FNV1A(key)
	if idx := slices.Index(m.key, hash); idx > -1 {
		m.key = slices.Delete(m.key, idx, idx+1)
		m.value = slices.Delete(m.value, idx, idx+1)
	}
	return m
}

/*
+-------+-------+-------+-------+-------+-------+
| len(8)|        key (64)       |    value      |  ...
+-------+-------+-------+-------+-------+-------+
*/

// Encode 编码 将字典编码为字节切片。
func MetaDictEncode(m MetaDict[string]) []byte {
	size := m.Len()
	if size == 0 {
		return nil
	}
	n := 0
	for i := 0; i < size; i++ {
		n += 1 + 8 + len(m.value[i])
	}
	buf := make([]byte, n)
	idx := 0
	for i := 0; i < size; i++ {
		buf[idx] = byte(1 + 8 + len(m.value[i]))
		idx++
		binary.LittleEndian.PutUint64(buf[idx:], m.key[i])
		idx += 8
		copy(buf[idx:], StringToBytes(m.value[i]))
		idx += len(m.value[i])
	}
	return buf
}

// Decode 解码 将字节切片解码为字典。
func MetaDictDecode(data []byte) (m MetaDict[string]) {
	idx := 0
	for len(data) > idx {
		n := int(data[idx])
		m.key = append(m.key, binary.LittleEndian.Uint64(data[idx+1:idx+1+8]))
		m.value = append(m.value, string(data[idx+1+8:idx+n]))
		idx += n
	}
	return
}
