package library

import (
	"errors"
	"fmt"
	"testing"

	"github.com/manifold/tractor/pkg/manifold/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testComponent struct {
	Foo string
}

func (c *testComponent) Echo(args ...string) []string {
	return args
}

func (c *testComponent) Err(msg string) (interface{}, error) {
	return nil, errors.New(msg)
}

func TestComponent(t *testing.T) {
	t.Run("GetSetFields", func(t *testing.T) {
		obj := &testComponent{Foo: "foo"}
		com := newComponent("test", obj, "")
		v, _, _ := com.GetField("Foo")
		assert.Equal(t, "foo", v)
		com.SetField("Foo", "bar")
		assert.Equal(t, "bar", obj.Foo)
	})
	t.Run("CallMethod", func(t *testing.T) {
		obj := &testComponent{Foo: "foo"}
		com := newComponent("test", obj, "")

		var echoRet []string
		noerr := com.CallMethod("Echo", []interface{}{"1", "2", "3"}, &echoRet)
		assert.Nil(t, noerr)
		assert.Len(t, echoRet, 3)

		err := com.CallMethod("Err", []interface{}{"error"}, nil)
		assert.Error(t, err)
		assert.Equal(t, "error", err.Error())
	})
}

type observingComponent struct {
	Name               string
	countEnabled       int
	countDisabled      int
	countSiblingReload map[interface{}]int
}

func (c *observingComponent) ComponentEnable() {
	c.countEnabled++
}

func (c *observingComponent) ComponentDisable() {
	c.countDisabled++
}

func (c *observingComponent) SiblingComponentReload(value interface{}) {
	if c.countSiblingReload == nil {
		c.countSiblingReload = make(map[interface{}]int)
	}
	name := fmt.Sprintf("%T", value)
	c.countSiblingReload[name]++
}

func (c *observingComponent) SiblingComponentReloads(name string) int {
	if c.countSiblingReload == nil {
		return 0
	}
	return c.countSiblingReload[name]
}

func TestComponentEnabledNotify(t *testing.T) {
	obj := object.New("root")
	refresher := &observingComponent{Name: "refresher"}
	listener := &observingComponent{Name: "listener"}
	otherCom := newComponent("other", &testComponent{Foo: "foo"}, "")
	refresherCom := newComponent("refresher", refresher, "")
	listenerCom := newComponent("listener", listener, "")
	obj.AppendComponent(otherCom)
	obj.AppendComponent(refresherCom)
	obj.AppendComponent(listenerCom)

	require.Nil(t, otherCom.Reload())
	require.Nil(t, refresherCom.Reload())
	require.Nil(t, listenerCom.Reload())

	assert.Equal(t, 1, refresher.countEnabled)
	assert.Equal(t, 0, refresher.countDisabled)
	assert.Equal(t, 1, refresher.SiblingComponentReloads("*library.observingComponent"))
	assert.Equal(t, 1, listener.countEnabled)
	assert.Equal(t, 0, listener.countDisabled)
	assert.Equal(t, 0, listener.SiblingComponentReloads("*library.observingComponent"))

	refresherCom.Reload()

	assert.Equal(t, 2, refresher.countEnabled)
	assert.Equal(t, 1, refresher.countDisabled)
	assert.Equal(t, 1, refresher.SiblingComponentReloads("*library.observingComponent"))
	assert.Equal(t, 1, listener.countEnabled)
	assert.Equal(t, 0, listener.countDisabled)
	assert.Equal(t, 1, listener.SiblingComponentReloads("*library.observingComponent"))

	for name, count := range listener.countSiblingReload {
		t.Logf("listener %q reloads: %d", name, count)
	}
}
