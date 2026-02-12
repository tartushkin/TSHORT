package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tStruct struct{ Value string }

func (m *tStruct) Reset() {
	m.Value = "reset"
}

var _ Resetter = (*tStruct)(nil) // проверка реализации

func TestPool(t *testing.T) {
	p := New[*tStruct]()
	obj := &tStruct{Value: "original"}

	p.Put(obj) // obj.Reset() → "reset"

	reused := p.Get()
	assert.Equal(t, "reset", reused.Value)
}
