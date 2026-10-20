package client

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ovn-kubernetes/libovsdb/ovsdb"
	"github.com/ovn-kubernetes/libovsdb/ovsdb/serverdb"
	"github.com/stretchr/testify/require"
)

func setDatabaseModel(t *testing.T, cli Client, row *serverdb.Database, dbModel string) {
	t.Helper()
	row.Model = dbModel
	ops, err := cli.Where(row).Update(row, &row.Model)
	require.NoError(t, err)
	reply, err := cli.Transact(t.Context(), ops...)
	require.NoError(t, err)
	_, err = ovsdb.CheckOperationResults(reply, ops)
	require.NoError(t, err)
}

func TestLeaderMonitorFollowsDatabaseModel(t *testing.T) {
	for _, initialModel := range []string{serverdb.DatabaseModelStandalone, serverdb.DatabaseModelRelay, serverdb.DatabaseModelClustered} {
		t.Run(initialModel, func(t *testing.T) {
			var connects, disconnects int32
			admin, row, endpoint := newClientServerPair(t, &connects, &disconnects, true)
			setDatabaseModel(t, admin, row, initialModel)
			c, err := newOVSDBClient(defDB, WithEndpoint(endpoint), WithLeaderOnly(true),
				WithReconnect(time.Second, backoff.NewConstantBackOff(time.Millisecond)))
			require.NoError(t, err)
			t.Cleanup(c.Close)
			require.NoError(t, c.Connect(t.Context()))

			checkMonitor := func(want bool) {
				t.Helper()
				c.rpcMutex.RLock()
				defer c.rpcMutex.RUnlock()
				require.Equal(t, want, c.leaderMonitorCookie != nil)
				db := c.databases[serverDB]
				db.monitorsMutex.Lock()
				defer db.monitorsMutex.Unlock()
				if want {
					require.Len(t, db.monitors, 1)
				} else {
					require.Empty(t, db.monitors)
				}
			}
			checkMonitor(initialModel == serverdb.DatabaseModelClustered)

			// A stable LB address can select a different model after reconnect.
			// Repeat the transition to ensure the watcher is neither lost nor duplicated.
			for _, nextModel := range []string{serverdb.DatabaseModelStandalone, serverdb.DatabaseModelClustered, serverdb.DatabaseModelRelay, serverdb.DatabaseModelClustered} {
				setDatabaseModel(t, admin, row, nextModel)
				before := atomic.LoadInt32(&connects)
				c.Disconnect()
				require.Eventually(t, func() bool {
					return atomic.LoadInt32(&connects) > before && c.Connected()
				}, 3*time.Second, time.Millisecond)
				checkMonitor(nextModel == serverdb.DatabaseModelClustered)
				require.NoError(t, c.Echo(t.Context()))
			}

			// The watcher installed after the model change must still reconnect.
			before := atomic.LoadInt32(&connects)
			setLeader(t, admin, row, false)
			require.Eventually(t, func() bool {
				return atomic.LoadInt32(&connects) > before
			}, 3*time.Second, time.Millisecond)
			setLeader(t, admin, row, true)
			require.Eventually(t, c.Connected, 3*time.Second, time.Millisecond)
			checkMonitor(true)

			// Explicitly connecting after Close must recreate the internal monitor.
			c.Close()
			select {
			case <-c.DisconnectNotify():
			case <-time.After(3 * time.Second):
				t.Fatal("close did not finish")
			}
			require.NoError(t, c.Connect(t.Context()))
			checkMonitor(true)
			require.NoError(t, c.Echo(t.Context()))
		})
	}
}

func TestLeaderOnlyRejectsSingleFollowerEndpoint(t *testing.T) {
	var connects, disconnects int32
	_, _, endpoint := newClientServerPair(t, &connects, &disconnects, false)
	c, err := newOVSDBClient(defDB, WithEndpoint(endpoint), WithLeaderOnly(true))
	require.NoError(t, err)
	t.Cleanup(c.Close)
	require.ErrorContains(t, c.Connect(t.Context()), "endpoint is not leader")
}
