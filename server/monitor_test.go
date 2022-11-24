package server

import (
	"testing"

	"github.com/kubeovn/libovsdb/ovsdb"
	"github.com/stretchr/testify/assert"
)

func TestMonitorFilter(t *testing.T) {
	monitor := monitor{
		request: map[string]*ovsdb.MonitorRequest{
			"Bridge": {
				Columns: []string{"name"},
				Select:  ovsdb.NewDefaultMonitorSelect(),
			},
		},
	}
	bridgeRow := ovsdb.Row{
		"_uuid": "foo",
		"name":  "bar",
	}
	bridgeRowWithIDs := ovsdb.Row{
		"_uuid":        "foo",
		"name":         "bar",
		"external_ids": map[string]string{"foo": "bar"},
	}
	tests := []struct {
		name     string
		update   ovsdb.TableUpdates2
		expected ovsdb.TableUpdates2
	}{
		{
			"not filtered",
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
		},
		{
			"removed table",
			ovsdb.TableUpdates2{
				"Open_vSwitch": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
			ovsdb.TableUpdates2{},
		},
		{
			"removed column",
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRowWithIDs,
					},
				},
			},
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor.filter2(tt.update)
			assert.Equal(t, tt.expected, tt.update)
		})
	}
}

func TestMonitorFilter2(t *testing.T) {
	monitor := monitor{
		request: map[string]*ovsdb.MonitorRequest{
			"Bridge": {
				Columns: []string{"name"},
				Select:  ovsdb.NewDefaultMonitorSelect(),
			},
		},
	}
	bridgeRow := ovsdb.Row{
		"name": "bar",
	}
	bridgeRowWithIDs := ovsdb.Row{
		"name":         "bar",
		"external_ids": map[string]string{"foo": "bar"},
	}
	tests := []struct {
		name     string
		update   ovsdb.TableUpdates2
		expected ovsdb.TableUpdates2
	}{
		{
			"not filtered",
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
		},
		{
			"removed table",
			ovsdb.TableUpdates2{
				"Open_vSwitch": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
			ovsdb.TableUpdates2{},
		},
		{
			"removed column",
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRowWithIDs,
					},
				},
			},
			ovsdb.TableUpdates2{
				"Bridge": ovsdb.TableUpdate2{
					"foo": &ovsdb.RowUpdate2{
						Insert: &bridgeRow,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor.filter2(tt.update)
			assert.Equal(t, tt.expected, tt.update)
		})
	}
}
