package utils

import (
	"strings"
	"testing"
)

func TestThree(t *testing.T) {
	if !strings.EqualFold(Three(true, "Yes", "No"), "Yes") {
		t.Error("不为 Yes")
	}
	if !strings.EqualFold(Three(false, "Yes", "No"), "No") {
		t.Error("不为 No")
	}
}
