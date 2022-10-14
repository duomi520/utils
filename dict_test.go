package utils

import (
	"testing"
)

var testKey []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var testValue []string = []string{"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9"}

func TestDict(t *testing.T) {
	var d MetaDict
	for i := 0; i < len(testKey); i++ {
		d.Set(testKey[i], testValue[i])
	}
	if d.Len() != 10 {
		t.Error("len != 10")
	}
	for i := 0; i < len(testKey); i++ {
		value, _ := d.Get(testKey[i])
		if value != testValue[i] {
			t.Error(value, testValue[i])
		}
	}
	d.Del("5")
	if d.Len() != 9 {
		t.Error("len != 9")
	}
	d.Del("0")
	if d.Len() != 8 {
		t.Error("len != 8")
	}
	d.Del("9")
	if d.Len() != 7 {
		t.Error("len != 7")
	}
}

func TestDictCode(t *testing.T) {
	var d, m MetaDict
	for i := 0; i < len(testKey); i++ {
		d.Set(testKey[i], testValue[i])
	}
	buf := make([]byte, 256)
	n, err := d.Encode(buf)
	if err != nil {
		t.Error(err)
	}
	m.Decode(buf[:n])
	if m.Len() != 10 {
		t.Error(m)
	}
	buf = make([]byte, 25)
	_, err = d.Encode(buf)
	if err == nil {
		t.Error("错误")
	}
}

func BenchmarkTestDict(b *testing.B) {
	var d MetaDict
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			d.Set(testKey[j], testValue[j])
			d.Get(testKey[j])
		}
	}
}
func BenchmarkTestMap(b *testing.B) {
	m := make(map[string]string, 16)
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			m[testKey[j]] = testValue[j]
			_ = m[testKey[j]]
		}
	}
}
