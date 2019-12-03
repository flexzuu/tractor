package manifold

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	sys := New()
	n1 := NewObject("n1")
	n2 := NewObject("n2")
	n3 := NewObject("n3")
	n2.AppendChild(n3)
	n1.AppendChild(n2)
	sys.AppendChild(n1)

	assert.Equal(t, sys, sys.System())
	assert.Equal(t, sys, n1.System())
	assert.Equal(t, sys, n2.System())
	assert.Equal(t, sys, n3.System())
}
