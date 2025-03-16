package utils

import (
	"bytes"
	"strings"
	"testing"
)

var testKey []string = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var testValue []string = []string{"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9"}

func TestDict(t *testing.T) {
	var d MetaDict[string]
	for i := range testKey {
		d = d.Set(testKey[i], testValue[i])
	}
	if d.Len() != 10 {
		t.Error("len != 10")
	}
	for i := range testKey {
		value, _ := d.Get(testKey[i])
		if value != testValue[i] {
			t.Error(value, testValue[i])
		}
	}
	d = d.Del("5")
	if d.Len() != 9 {
		t.Error("len != 9")
	}
	d = d.Del("0")
	if d.Len() != 8 {
		t.Error("len != 8")
	}
	d = d.Del("9")
	if d.Len() != 7 {
		t.Error("len != 7")
	}
}

func TestDictCode(t *testing.T) {
	var d MetaDict[string]
	for i := range testKey {
		d = d.Set(testKey[i], testValue[i])
	}
	buf := MetaDictEncode(d)
	m1 := MetaDictDecode(buf)
	if m1.Len() != 10 {
		t.Error(m1)
	}
	v, ok := m1.Get("7")
	if !strings.EqualFold(v, "a7") {
		t.Error(v, ok)
	}
	buffer := bytes.Buffer{}
	err := MetaDictEncoder(d, &buffer)
	if err != nil {
		t.Fatal(err.Error())
	}
	m2 := MetaDictDecode(buffer.Bytes())
	if m2.Len() != 10 {
		t.Error(m2)
	}
	v, ok = m2.Get("5")
	if !strings.EqualFold(v, "a5") {
		t.Error(v, ok)
	}
}

func BenchmarkTestDict(b *testing.B) {
	var d MetaDict[string]
	for i := 0; i < b.N; i++ {
		for j := range 10 {
			d = d.Set(testKey[j], testValue[j])
			d.Get(testKey[j])
		}
	}
}
func BenchmarkTestMap(b *testing.B) {
	m := make(map[string]string, 16)
	for i := 0; i < b.N; i++ {
		for j := range 10 {
			m[testKey[j]] = testValue[j]
			_ = m[testKey[j]]
		}
	}
}
