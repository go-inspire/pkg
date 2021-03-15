package log

import (
	"testing"
)

func TestCrit(t *testing.T) {
	Debugf("test %d", 1)
	Infof("test")
}
