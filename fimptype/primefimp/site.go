package primefimp

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Site object
type Site struct {
	ID        int          `json:"id,omitempty"`
	Devices   []Device     `json:"device,omitempty"`
	Things    []Thing      `json:"thing,omitempty"`
	Rooms     []Room       `json:"room,omitempty"`
	House     *House       `json:"house,omitempty"`
	Hub       *Hub         `json:"hub,omitempty"`
	Areas     []Area       `json:"area,omitempty"`
	Shortcuts []Shortcut   `json:"shortcut,omitempty"`
	Services  VincServices `json:"service,omitempty"`
	Modes     []Mode       `json:"mode,omitempty"`
	Timers    []Timer      `json:"timer,omitempty"`
	Problem   bool         `json:"problem,omitempty"`
}

// SiteFromResponse Creates a Site object from given response
func SiteFromResponse(resp *Response) *Site {
	site := Site{Devices: resp.GetDevices(), Things: resp.GetThings(), Rooms: resp.GetRooms(), Areas: resp.GetAreas(), House: resp.GetHouse(),
		Shortcuts: resp.GetShortcuts(), Modes: resp.GetModes(), Timers: resp.GetTimers(), Services: resp.GetVincServices()}
	return &site
}

// AddDevice Adds device
func (s *Site) AddDevice(d *Device) {
	if s.FindIndex(ComponentDevice, d.ID) == -1 {
		s.Devices = append(s.Devices, *d)
	}
}

// AddRoom Adds room
func (s *Site) AddRoom(r *Room) {
	if s.FindIndex(ComponentRoom, r.ID) == -1 {
		s.Rooms = append(s.Rooms, *r)
	}
}

// AddArea Adds area
func (s *Site) AddArea(a *Area) {
	if s.FindIndex(ComponentArea, a.ID) == -1 {
		s.Areas = append(s.Areas, *a)
	}
}

// AddTimer Adds timer
func (s *Site) AddTimer(ti *Timer) {
	if s.FindIndex(ComponentTimer, ti.ID) == -1 {
		s.Timers = append(s.Timers, *ti)
	}
}

// AddThing Adds thing
func (s *Site) AddThing(th *Thing) {
	if s.FindIndex(ComponentThing, th.ID) == -1 {
		s.Things = append(s.Things, *th)
	}
}

// AddShortcut Adds shortcut
func (s *Site) AddShortcut(sh *Shortcut) {
	if s.FindIndex(ComponentShortcut, sh.ID) == -1 {
		s.Shortcuts = append(s.Shortcuts, *sh)
	}
}

// FindIndex Finds the index of a component with given ID
func (s *Site) FindIndex(comp string, id int) int {
	switch comp {
	case ComponentArea:
		for k, v := range s.Areas {
			if id == v.ID {
				return k
			}
		}
	case ComponentDevice:
		for k, v := range s.Devices {
			if id == v.ID {
				return k
			}
		}
	case ComponentRoom:
		for k, v := range s.Rooms {
			if id == v.ID {
				return k
			}
		}
	case ComponentTimer:
		for k, v := range s.Timers {
			if id == v.ID {
				return k
			}
		}
	case ComponentThing:
		for k, v := range s.Things {
			if id == v.ID {
				return k
			}
		}
	case ComponentShortcut:
		for k, v := range s.Shortcuts {
			if id == v.ID {
				return k
			}
		}
	default:
		log.Error("Component does not support findindex.")
	}
	return -1
}

// RemoveWithID Removes the component with given ID.
func (s *Site) RemoveWithID(comp string, id int) error {
	switch comp {
	case ComponentArea:
		idx := -1
		for k, v := range s.Areas {
			if id == v.ID {
				idx = k
			}
		}
		s.Areas[idx] = s.Areas[len(s.Areas)-1]
		s.Areas = s.Areas[:len(s.Areas)-1]
	case ComponentDevice:
		idx := -1
		for k, v := range s.Devices {
			if id == v.ID {
				idx = k
			}
		}
		s.Devices[idx] = s.Devices[len(s.Devices)-1]
		s.Devices = s.Devices[:len(s.Devices)-1]
	case ComponentRoom:
		idx := -1
		for k, v := range s.Rooms {
			if id == v.ID {
				idx = k
			}
		}
		s.Rooms[idx] = s.Rooms[len(s.Rooms)-1]
		s.Rooms = s.Rooms[:len(s.Rooms)-1]
	case ComponentTimer:
		idx := -1
		for k, v := range s.Timers {
			if id == v.ID {
				idx = k
			}
		}
		s.Timers[idx] = s.Timers[len(s.Timers)-1]
		s.Timers = s.Timers[:len(s.Timers)-1]
	case ComponentThing:
		idx := -1
		for k, v := range s.Things {
			if id == v.ID {
				idx = k
			}
		}
		s.Things[idx] = s.Things[len(s.Things)-1]
		s.Things = s.Things[:len(s.Things)-1]
	case ComponentShortcut:
		idx := -1
		for k, v := range s.Shortcuts {
			if id == v.ID {
				idx = k
			}
		}
		s.Shortcuts[idx] = s.Shortcuts[len(s.Shortcuts)-1]
		s.Shortcuts = s.Shortcuts[:len(s.Shortcuts)-1]
	default:
		return fmt.Errorf("RemoveWithID: %s does not support removal", comp)

	}
	return fmt.Errorf("RemoveWithID: %s with ID:%d not found", comp, id)
}
