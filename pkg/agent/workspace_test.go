package agent

import (
	"testing"

	"github.com/manifold/tractor/pkg/workspace/supervisor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspace(t *testing.T) {
	ag, teardown := setup(t, "test1", "test2", "test3")
	defer teardown()

	t.Run("stop/start", func(t *testing.T) {
		status, ws := setupWorkspace(t, ag, "test1")
		assert.Equal(t, supervisor.StatusAvailable, <-status)

		assert.Nil(t, ws.Stop())
		assert.Equal(t, supervisor.StatusUnavailable, <-status)

		assert.Nil(t, ws.Start())
		assert.Equal(t, supervisor.StatusAvailable, <-status)
	})

	// t.Run("connect/stop", func(t *testing.T) {
	// 	ws := ag.Workspace("test3")
	// 	require.NotNil(t, ws)
	// 	assert.Equal(t, StatusAvailable, ws.Status)

	// 	connCh := readWorkspace(t, ws.Connect)
	// 	time.Sleep(time.Second)
	// 	assert.Equal(t, StatusAvailable, ws.Status)

	// 	ws.Stop()
	// 	assert.Equal(t, StatusUnavailable, ws.Status)

	// 	connOut := strings.TrimSpace(string(<-connCh))
	// 	assert.True(t, strings.HasPrefix(connOut, "pid "))
	// })

}

func setupWorkspace(t *testing.T, ag *Agent, name string) (chan supervisor.Status, *Workspace) {
	status := make(chan supervisor.Status, 3)
	ws := ag.Workspace(name)
	require.NotNil(t, ws)
	//ws.SetDaemonCmd("cat")
	ws.Observe(func(_ *supervisor.Supervisor, newStatus supervisor.Status) {
		status <- newStatus
	})
	require.NoError(t, ws.StartDaemon())
	return status, ws
}
