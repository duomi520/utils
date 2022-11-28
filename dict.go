package utils

import (
	"errors"
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

//GetAll
func (d *MetaDict) GetAll() ([]string, []string) {
	return d.key, d.value
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
func (d *MetaDict) Encode(src []byte) (int, error) {
	index := 0
	for i := 0; i < d.lenght; i++ {
		need := 2 + len(d.key[i]) + len(d.value[i])
		if len(src) < (index + need) {
			return 0, errors.New("MetaDict.Encode：[]byte is too short")
		}
		src[index] = byte(len(d.key[i]))
		index++
		copy(src[index:], StringToBytes(d.key[i]))
		index = index + len(d.key[i])
		src[index] = byte(len(d.value[i]))
		index++
		copy(src[index:], StringToBytes(d.value[i]))
		index = index + len(d.value[i])
	}
	return index, nil
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
