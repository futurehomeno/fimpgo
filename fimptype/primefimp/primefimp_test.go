package primefimp

import (
	"encoding/json"
	"testing"

	"github.com/futurehomeno/fimpgo"
	"github.com/futurehomeno/fimpgo/fimptype"
	"github.com/futurehomeno/fimpgo/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMsg(t *testing.T, intf string, val any) *fimpgo.FimpMessage {
	t.Helper()
	b, err := json.Marshal(val)
	require.NoError(t, err)
	msg := fimpgo.NewMessage(intf, "vinculum", fimptype.VTypeObject, nil, nil, nil, nil)
	msg.ValueObj = b
	return msg
}

func newResponse(t *testing.T, param map[string]any) *Response {
	t.Helper()
	resp, err := FimpToResponse(newMsg(t, "evt.pd7.response", map[string]any{"cmd": "get", "success": true, "param": param}))
	require.NoError(t, err)
	return resp
}

func fullResponse(t *testing.T) *Response {
	t.Helper()
	return newResponse(t, map[string]any{
		"device":   []Device{{ID: 1}, {ID: 2, Fimp: Fimp{Address: "a"}, Service: map[string]Service{"s": {Addr: "/rt:dev/ad:1"}}}},
		"thing":    []Thing{{ID: 3}},
		"room":     []Room{{ID: 4}},
		"area":     []Area{{ID: 5}},
		"shortcut": []Shortcut{{ID: 6}},
		"timer":    []Timer{{ID: 7}},
		"mode":     []Mode{{Id: "home"}},
		"house":    House{Mode: "home"},
		"service":  VincServices{FireAlarm: map[string]any{"a": 1}},
	})
}

func TestFimpToResponse(t *testing.T) {
	_, err := FimpToResponse(newMsg(t, "evt.other", map[string]any{}))
	require.Error(t, err)
	_, err = FimpToResponse(fimpgo.NewStringMessage("evt.pd7.response", "vinculum", "x", nil, nil, nil))
	require.Error(t, err)
}

func TestResponseGetters(t *testing.T) {
	resp := fullResponse(t)
	assert.Len(t, resp.GetDevices(), 2)
	assert.Len(t, resp.GetThings(), 1)
	assert.Len(t, resp.GetRooms(), 1)
	assert.Len(t, resp.GetAreas(), 1)
	assert.Len(t, resp.GetShortcuts(), 1)
	assert.Len(t, resp.GetTimers(), 1)
	assert.Equal(t, "home", resp.GetModes()[0].Id)
	assert.Equal(t, "home", resp.GetHouse().Mode)
	assert.Equal(t, 1.0, resp.GetVincServices().FireAlarm["a"])

	empty := newResponse(t, map[string]any{})
	assert.Nil(t, empty.GetDevices())
	assert.Nil(t, empty.GetThings())
	assert.Nil(t, empty.GetRooms())
	assert.Nil(t, empty.GetAreas())
	assert.Nil(t, empty.GetShortcuts())
	assert.Nil(t, empty.GetModes())
	assert.Nil(t, empty.GetTimers())
	assert.Nil(t, empty.GetHouse())
	assert.Empty(t, empty.GetVincServices())
	st, err := empty.GetState()
	require.NoError(t, err)
	assert.Empty(t, st.Devices)

	bad := newResponse(t, map[string]any{"device": "x", "thing": "x", "room": "x", "area": "x", "shortcut": "x",
		"mode": "x", "timer": "x", "house": "x", "service": "x", "state": "x"})
	assert.Nil(t, bad.GetDevices())
	assert.Nil(t, bad.GetThings())
	assert.Nil(t, bad.GetRooms())
	assert.Nil(t, bad.GetAreas())
	assert.Nil(t, bad.GetShortcuts())
	assert.Nil(t, bad.GetModes())
	assert.Nil(t, bad.GetTimers())
	assert.Nil(t, bad.GetHouse())
	assert.Empty(t, bad.GetVincServices())
	_, err = bad.GetState()
	require.Error(t, err)
}

func TestSiteFromResponseAndLookups(t *testing.T) {
	site := SiteFromResponse(fullResponse(t))
	assert.Equal(t, 3, site.GetThingById(3).ID)
	assert.Nil(t, site.GetThingById(99))
	assert.Equal(t, 4, site.GetRoomById(4).ID)
	assert.Nil(t, site.GetRoomById(99))
	assert.Equal(t, 5, site.GetAreaById(5).ID)
	assert.Nil(t, site.GetAreaById(99))
	assert.Equal(t, 2, site.GetDeviceById(2).ID)
	assert.Nil(t, site.GetDeviceById(99))
	assert.Equal(t, 2, site.GetDeviceByServiceAddress("/rt:dev/ad:1").ID)
	assert.Nil(t, site.GetDeviceByServiceAddress("nope"))
	assert.Equal(t, "/rt:dev/ad:1", site.GetServiceByAddress("/rt:dev/ad:1").Addr)
	assert.Nil(t, site.GetServiceByAddress("nope"))
}

func TestSiteAddFindUpdateRemove(t *testing.T) {
	adders := map[string]func(s *Site){
		ComponentDevice:   func(s *Site) { s.AddDevice(&Device{ID: 1}) },
		ComponentRoom:     func(s *Site) { s.AddRoom(&Room{ID: 1}) },
		ComponentArea:     func(s *Site) { s.AddArea(&Area{ID: 1}) },
		ComponentTimer:    func(s *Site) { s.AddTimer(&Timer{ID: 1}) },
		ComponentThing:    func(s *Site) { s.AddThing(&Thing{ID: 1}) },
		ComponentShortcut: func(s *Site) { s.AddShortcut(&Shortcut{ID: 1}) },
	}
	updaters := map[string]func(s *Site){
		ComponentDevice:   func(s *Site) { s.UpdateDevice(&Device{ID: 1, Model: "m"}) },
		ComponentRoom:     func(s *Site) { s.UpdateRoom(&Room{ID: 1, Alias: "m"}) },
		ComponentArea:     func(s *Site) { s.UpdateArea(&Area{ID: 1, Name: "m"}) },
		ComponentTimer:    func(s *Site) { s.UpdateTimer(&Timer{ID: 1, Enabled: true}) },
		ComponentThing:    func(s *Site) { s.UpdateThing(&Thing{ID: 1, Name: "m"}) },
		ComponentShortcut: func(s *Site) { s.UpdateShortcut(&Shortcut{ID: 1, Client: Client{Name: new(string)}}) },
	}
	for comp, add := range adders {
		t.Run(comp, func(t *testing.T) {
			s := &Site{}
			assert.Equal(t, -1, s.FindIndex(comp, 1))
			updaters[comp](s)
			assert.Equal(t, 0, s.FindIndex(comp, 1))
			s2 := &Site{}
			add(s2)
			add(s2)
			assert.Equal(t, 0, s2.FindIndex(comp, 1))
			assert.Equal(t, -1, s2.FindIndex(comp, 2))
			updaters[comp](s2)
			require.NoError(t, s2.RemoveWithID(comp, 1))
			assert.Equal(t, -1, s2.FindIndex(comp, 1))
			require.Error(t, s2.RemoveWithID(comp, 1))
		})
	}
	s := &Site{Devices: []Device{{ID: 1}, {ID: 2}, {ID: 3}}}
	s.UpdateDevice(&Device{ID: 2, Model: "x"})
	assert.Equal(t, "x", s.Devices[1].Model)
	require.NoError(t, s.RemoveWithID(ComponentDevice, 1))
	assert.Len(t, s.Devices, 2)
	assert.NotNil(t, s.GetDeviceById(3))
	require.Error(t, s.RemoveWithID(ComponentHub, 1))
	assert.Equal(t, -1, s.FindIndex(ComponentHub, 1))
}

func newNotify(t *testing.T, cmd, comp string, param, changes any) *Notify {
	t.Helper()
	pb, err := json.Marshal(param)
	require.NoError(t, err)
	cb, err := json.Marshal(changes)
	require.NoError(t, err)
	return &Notify{Cmd: cmd, Component: comp, ParamRaw: pb, ChangesRaw: cb}
}

func TestFimpToNotify(t *testing.T) {
	_, err := FimpToNotify(&fimpgo.Message{Payload: newMsg(t, "evt.other", map[string]any{})})
	require.Error(t, err)
	_, err = FimpToNotify(&fimpgo.Message{Payload: fimpgo.NewStringMessage("evt.pd7.notify", "vinculum", "x", nil, nil, nil)})
	require.Error(t, err)

	n, err := FimpToNotify(&fimpgo.Message{Payload: newMsg(t, "evt.pd7.notify", map[string]any{"cmd": "edit", "component": "room", "param": Room{ID: 8}})})
	require.NoError(t, err)
	assert.Equal(t, "edit", n.Cmd)
	assert.Equal(t, 8, n.GetRoom().ID)
}

func TestNotifyGetters(t *testing.T) {
	assert.Equal(t, 1, newNotify(t, CmdAdd, ComponentDevice, Device{ID: 1}, nil).GetDevice().ID)
	assert.Equal(t, 1, newNotify(t, CmdAdd, ComponentThing, Thing{ID: 1}, nil).GetThing().ID)
	assert.Equal(t, 1, newNotify(t, CmdAdd, ComponentRoom, Room{ID: 1}, nil).GetRoom().ID)
	assert.Equal(t, 1, newNotify(t, CmdAdd, ComponentArea, Area{ID: 1}, nil).GetArea().ID)
	assert.Equal(t, "m", newNotify(t, CmdAdd, ComponentHouse, House{Mode: "m"}, nil).GetHouse().Mode)
	assert.Equal(t, 1, newNotify(t, CmdAdd, ComponentShortcut, Shortcut{ID: 1}, nil).GetShortcut().ID)
	assert.Equal(t, 1, newNotify(t, CmdAdd, ComponentTimer, Timer{ID: 1}, nil).GetTimer().ID)
	assert.Equal(t, "away", newNotify(t, CmdSet, ComponentHub, Hub{Mode: HubMode{Current: "away"}}, nil).GetHub().Mode.Current)
	assert.Equal(t, 7, newNotify(t, CmdDelete, ComponentRoom, nil, DeleteChange{ID: 7}).GetDeleteChange().ID)
	assert.Equal(t, "b", newNotify(t, CmdSet, ComponentHouse, ModeChange{Current: "b", Prev: "a"}, nil).GetModeChange().Current)

	wrong := newNotify(t, CmdAdd, ComponentHub, Hub{}, nil)
	assert.Nil(t, wrong.GetDevice())
	assert.Nil(t, wrong.GetThing())
	assert.Nil(t, wrong.GetRoom())
	assert.Nil(t, wrong.GetArea())
	assert.Nil(t, wrong.GetHouse())
	assert.Nil(t, wrong.GetShortcut())
	assert.Nil(t, wrong.GetTimer())

	for comp, get := range map[string]func(n *Notify) any{
		ComponentDevice:   func(n *Notify) any { return n.GetDevice() },
		ComponentThing:    func(n *Notify) any { return n.GetThing() },
		ComponentRoom:     func(n *Notify) any { return n.GetRoom() },
		ComponentArea:     func(n *Notify) any { return n.GetArea() },
		ComponentHouse:    func(n *Notify) any { return n.GetHouse() },
		ComponentShortcut: func(n *Notify) any { return n.GetShortcut() },
		ComponentTimer:    func(n *Notify) any { return n.GetTimer() },
		ComponentHub:      func(n *Notify) any { return n.GetHub() },
	} {
		bad := &Notify{Component: comp, ParamRaw: json.RawMessage(`"x"`)}
		assert.Nil(t, get(bad), comp)
	}
	broken := &Notify{}
	assert.Nil(t, broken.GetDeleteChange())
	assert.Nil(t, broken.GetModeChange())
}

func TestUpdateSite(t *testing.T) {
	api := NewApiClient("c", nil, false)
	api.EnableVincAppsSync(true)

	cases := []struct {
		name    string
		notif   *Notify
		wantErr bool
	}{
		{"add area", newNotify(t, CmdAdd, ComponentArea, Area{ID: 1}, nil), false},
		{"add device", newNotify(t, CmdAdd, ComponentDevice, Device{ID: 1}, nil), false},
		{"add room", newNotify(t, CmdAdd, ComponentRoom, Room{ID: 1}, nil), false},
		{"add thing", newNotify(t, CmdAdd, ComponentThing, Thing{ID: 1}, nil), false},
		{"add shortcut", newNotify(t, CmdAdd, ComponentShortcut, Shortcut{ID: 1}, nil), false},
		{"add timer", newNotify(t, CmdAdd, ComponentTimer, Timer{ID: 1}, nil), false},
		{"add unknown", newNotify(t, CmdAdd, ComponentHub, Hub{}, nil), true},
		{"edit area", newNotify(t, CmdEdit, ComponentArea, Area{ID: 1, Name: "n"}, nil), false},
		{"edit device", newNotify(t, CmdEdit, ComponentDevice, Device{ID: 1, Model: "n"}, nil), false},
		{"edit room", newNotify(t, CmdEdit, ComponentRoom, Room{ID: 1, Alias: "n"}, nil), false},
		{"edit thing", newNotify(t, CmdEdit, ComponentThing, Thing{ID: 1, Name: "n"}, nil), false},
		{"edit timer", newNotify(t, CmdEdit, ComponentTimer, Timer{ID: 1, Enabled: true}, nil), false},
		{"edit shortcut", newNotify(t, CmdEdit, ComponentShortcut, Shortcut{ID: 1}, nil), false},
		{"edit unknown", newNotify(t, CmdEdit, ComponentHub, Hub{}, nil), true},
		{"set room", newNotify(t, CmdSet, ComponentRoom, nil, nil), false},
		{"set hub", newNotify(t, CmdSet, ComponentHub, nil, nil), false},
		{"set device", newNotify(t, CmdSet, ComponentDevice, nil, nil), false},
		{"set unknown", newNotify(t, CmdSet, ComponentThing, nil, nil), true},
		{"unknown cmd", newNotify(t, "bogus", ComponentRoom, nil, nil), true},
		{"delete bad id type", &Notify{Cmd: CmdDelete, Component: ComponentRoom, Id: "1"}, true},
		{"delete missing", &Notify{Cmd: CmdDelete, Component: ComponentRoom, Id: 99.0}, true},
		{"delete float id", &Notify{Cmd: CmdDelete, Component: ComponentRoom, Id: 1.0}, false},
		{"delete int id", &Notify{Cmd: CmdDelete, Component: ComponentArea, Id: 1}, false},
	}
	for _, tc := range cases {
		err := api.UpdateSite(tc.notif)
		assert.Equal(t, tc.wantErr, err != nil, tc.name)
	}
	assert.Equal(t, "n", api.siteCache.Devices[0].Model)
	assert.Empty(t, api.siteCache.Rooms)

	api.EnableVincAppsSync(false)
	require.NoError(t, api.UpdateSite(newNotify(t, CmdAdd, ComponentShortcut, Shortcut{ID: 2}, nil)))
	assert.Len(t, api.siteCache.Shortcuts, 1)
	require.NoError(t, api.UpdateSite(newNotify(t, "bogus", ComponentHub, nil, nil)))
}

func TestApiClientCache(t *testing.T) {
	api := NewApiClient("c", nil, false, WithCloudService("svc"), WithGlobalPrefix("g"))
	assert.Equal(t, "svc", api.cloudService)
	assert.Equal(t, "g", api.globalPrefix)
	assert.True(t, api.IsCacheEmpty())
	assert.False(t, api.IsCacheEnabled())

	_, err := api.GetDevices(true)
	require.ErrorIs(t, err, utils.ErrEmptyCache)

	require.NoError(t, api.LoadVincResponseFromFile("testdata/site-info-response.json"))
	assert.False(t, api.IsCacheEmpty())
	assert.True(t, api.ValidateAndReloadSiteCache())
	_, err = api.GetDevices(true)
	require.ErrorIs(t, err, utils.ErrEmptyCache, "cache disabled")

	api.SetIsCacheEnabled(true)
	assert.True(t, api.IsCacheEnabled())
	api.siteCache.House = &House{Mode: "home"}
	api.siteCache.Shortcuts = []Shortcut{{ID: 1}}
	api.siteCache.Modes = []Mode{{Id: "m"}}
	api.siteCache.Timers = []Timer{{ID: 1}}

	for name, get := range map[string]func() (any, error){
		"devices":   func() (any, error) { return api.GetDevices(true) },
		"rooms":     func() (any, error) { return api.GetRooms(true) },
		"areas":     func() (any, error) { return api.GetAreas(true) },
		"things":    func() (any, error) { return api.GetThings(true) },
		"shortcuts": func() (any, error) { return api.GetShortcuts(true) },
		"modes":     func() (any, error) { return api.GetModes(true) },
		"mode":      func() (any, error) { return api.GetCurrentMode(true) },
		"timers":    func() (any, error) { return api.GetTimers(true) },
		"services":  func() (any, error) { return api.GetVincServices(true) },
		"site":      func() (any, error) { return api.GetSite(true) },
	} {
		got, err := get()
		require.NoError(t, err, name)
		assert.NotNil(t, got, name)
	}
	house, _ := api.GetCurrentMode(true)
	assert.Equal(t, "home", house.Mode)

	api.siteCache = Site{}
	assert.True(t, api.IsCacheEmpty())
}

func TestApiClientCacheDisabledErrors(t *testing.T) {
	api := NewApiClient("c", nil, false)
	for name, get := range map[string]func() error{
		"rooms":     func() error { _, err := api.GetRooms(true); return err },
		"areas":     func() error { _, err := api.GetAreas(true); return err },
		"things":    func() error { _, err := api.GetThings(true); return err },
		"shortcuts": func() error { _, err := api.GetShortcuts(true); return err },
		"modes":     func() error { _, err := api.GetModes(true); return err },
		"mode":      func() error { _, err := api.GetCurrentMode(true); return err },
		"timers":    func() error { _, err := api.GetTimers(true); return err },
		"services":  func() error { _, err := api.GetVincServices(true); return err },
	} {
		require.ErrorIs(t, get(), utils.ErrEmptyCache, name)
	}
}

func TestApiClientLoadFromFileErrors(t *testing.T) {
	api := NewApiClient("c", nil, false)
	require.Error(t, api.LoadVincResponseFromFile("testdata/missing.json"))
	require.Error(t, api.LoadVincResponseFromFile("testdata/mode.json"))
	require.Error(t, api.LoadVincResponseFromFile("primefimp_test.go"))
}

func TestApiClientChannels(t *testing.T) {
	api := NewApiClient("c", nil, false)
	ch := make(chan Notify, 1)
	api.RegisterChannel("a", ch)
	api.RegisterChannelWithFilter("b", ch, NotifyFilter{Cmd: CmdAdd, Component: ComponentRoom})
	assert.Len(t, api.notifySubChannels, 2)
	assert.Len(t, api.subFilters, 1)
	api.UnregisterChannel("b")
	assert.Len(t, api.notifySubChannels, 1)
	assert.Empty(t, api.subFilters)
}

func TestApiClientResponseAddress(t *testing.T) {
	api := NewApiClient("cid", nil, false)
	api.SetResponsePayloadType("j2")
	addr := api.responseAddress()
	assert.Equal(t, "pt:j2/mt:rsp/rt:app/rn:cid/ad:1", addr.Serialize())

	api = NewApiClient("cid", nil, false, WithCloudService("svc"))
	addr = api.responseAddress()
	assert.Equal(t, "pt:j1/mt:rsp/rt:cloud/rn:backend-service/ad:svc", addr.Serialize())
}
