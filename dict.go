package utils

import (
	"strings"
)

//MetaDict 非线程安全,key数量超过5个后，效率低于map
type MetaDict struct {
	lenght int
	key    []string
	value  []string
}

//Len
func (d *MetaDict) Len() int {
	return d.lenght
}

//Set
func (d *MetaDict) Set(key, value string) {
	for i := 0; i < d.lenght; i++ {
		if d.key[i] == key {
			d.value[i] = value
			return
		}
	}
	d.key = append(d.key, key)
	d.value = append(d.value, value)
	d.lenght++
}

//Get
func (d *MetaDict) Get(key string) (string, bool) {
	for i := 0; i < d.lenght; i++ {
		if d.key[i] == key {
			return d.value[i], true
		}
	}
	return "", false
}

//Del
func (d *MetaDict) Del(key string) {
	for i := 0; i < d.lenght; i++ {
		if strings.EqualFold(d.key[i], key) {
			if i < (d.lenght - 1) {
				copy(d.key[i:], d.key[i+1:])
				copy(d.value[i:], d.value[i+1:])
			}
			d.lenght--
			d.key = d.key[:d.lenght]
			d.value = d.value[:d.lenght]
			return
		}
	}
}

/*
+-------+-------+-------+-------+-------+-------+
| len(8)|      key      | len(8)|    value      |  ...
+-------+-------+-------+-------+-------+-------+
*/

//Encode 编码
func (d *MetaDict) Encode(b []byte) int {
	index := 0
	for i := 0; i < d.lenght; i++ {
		b[index] = byte(len(d.key[i]))
		index++
		copy(b[index:], StringToBytes(d.key[i]))
		index = index + len(d.key[i])
		b[index] = byte(len(d.value[i]))
		index++
		copy(b[index:], StringToBytes(d.value[i]))
		index = index + len(d.value[i])
	}
	return index
}

//Decode 解码
func (d *MetaDict) Decode(data []byte) {
	index, size := 0, 0
	for len(data) > index {
		size = int(data[index])
		index++
		d.key = append(d.key, string(data[index:index+size]))
		index = index + size
		size = int(data[index])
		index++
		d.value = append(d.value, string(data[index:index+size]))
		index = index + size
	}
	d.lenght = len(d.key)
}
