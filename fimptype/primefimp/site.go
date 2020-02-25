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
// returns the index in the arreay of corresponding component
// return -1 in case the component is not found
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
		idx := s.FindIndex(comp, id)
		if idx != -1 {
			s.Areas[idx] = s.Areas[len(s.Areas)-1]
			s.Areas = s.Areas[:len(s.Areas)-1]
		} else {
			return fmt.Errorf("RemoveWithID: Area ID:%d not found", id)
		}
	case ComponentDevice:
		idx := s.FindIndex(comp, id)
		if idx != -1 {
			s.Devices[idx] = s.Devices[len(s.Devices)-1]
			s.Devices = s.Devices[:len(s.Devices)-1]
		} else {
			return fmt.Errorf("RemoveWithID: Device ID:%d not found", id)
		}
	case ComponentRoom:
		idx := s.FindIndex(comp, id)
		if idx != -1 {
			s.Rooms[idx] = s.Rooms[len(s.Rooms)-1]
			s.Rooms = s.Rooms[:len(s.Rooms)-1]
		} else {
			return fmt.Errorf("RemoveWithID: Room ID:%d not found", id)
		}
	case ComponentTimer:
		idx := s.FindIndex(comp, id)
		if idx != -1 {
			s.Timers[idx] = s.Timers[len(s.Timers)-1]
			s.Timers = s.Timers[:len(s.Timers)-1]
		} else {
			return fmt.Errorf("RemoveWithID: Timer ID:%d not found", id)
		}
	case ComponentThing:
		idx := s.FindIndex(comp, id)
		if idx != -1 {
			s.Things[idx] = s.Things[len(s.Things)-1]
			s.Things = s.Things[:len(s.Things)-1]
		} else {
			return fmt.Errorf("RemoveWithID: Thing ID:%d not found", id)
		}
	case ComponentShortcut:
		idx := s.FindIndex(comp, id)
		if idx != -1 {
			s.Shortcuts[idx] = s.Shortcuts[len(s.Shortcuts)-1]
			s.Shortcuts = s.Shortcuts[:len(s.Shortcuts)-1]
		} else {
			return fmt.Errorf("RemoveWithID: Shortcut ID:%d not found", id)
		}
	default:
		return fmt.Errorf("RemoveWithID: %s does not support removal", comp)

	}
	return nil
}

// UpdateDevice Updates device
func (s *Site) UpdateDevice(d *Device) {
	idx := s.FindIndex(ComponentDevice, d.ID)
	if idx == -1 {
		// if component is not added before (somehow) add it.
		s.Devices = append(s.Devices, *d)
	} else {
		s.Devices[idx] = *d
	}
}

// UpdateArea Updates Area
func (s *Site) UpdateArea(d *Area) {
	idx := s.FindIndex(ComponentArea, d.ID)
	if idx == -1 {
		// if component is not added before (somehow) add it.
		s.Areas = append(s.Areas, *d)
	} else {
		s.Areas[idx] = *d
	}
}

// UpdateRoom Updates Room
func (s *Site) UpdateRoom(d *Room) {
	idx := s.FindIndex(ComponentRoom, d.ID)
	if idx == -1 {
		// if component is not added before (somehow) add it.
		s.Rooms = append(s.Rooms, *d)
	} else {
		s.Rooms[idx] = *d
	}
}

// UpdateThing Updates Thing
func (s *Site) UpdateThing(d *Thing) {
	idx := s.FindIndex(ComponentThing, d.ID)
	if idx == -1 {
		// if component is not added before (somehow) add it.
		s.Things = append(s.Things, *d)
	} else {
		s.Things[idx] = *d
	}
}

// UpdateTimer Updates Timer
func (s *Site) UpdateTimer(d *Timer) {
	idx := s.FindIndex(ComponentTimer, d.ID)
	if idx == -1 {
		// if component is not added before (somehow) add it.
		s.Timers = append(s.Timers, *d)
	} else {
	    s.Timers[idx] = *d
	}
}

// UpdateShortcut Updates Shortcut
func (s *Site) UpdateShortcut(d *Shortcut) {
	idx := s.FindIndex(ComponentShortcut, d.ID)
	if idx == -1 {
		// if component is not added before (somehow) add it.
		s.Shortcuts = append(s.Shortcuts, *d)
	} else {
		s.Shortcuts[idx] = *d
	}
}

func (s *Site) GetRoomById(ID int)*Room {
	for i := range s.Rooms {
		if s.Rooms[i].ID == ID {
			return &s.Rooms[i]
		}
	}
	return nil
}

func (s *Site) GetAreaById(ID int)*Area {
	for i := range s.Areas {
		if s.Areas[i].ID == ID {
			return &s.Areas[i]
		}
	}
	return nil
}

func (s *Site) GetThingById(ID int)*Thing {
	for i := range s.Things {
		if s.Things[i].ID == ID {
			return &s.Things[i]
		}
	}
	return nil
}

func (s *Site) GetDeviceById(ID int)*Device {
	for i := range s.Devices {
		if s.Devices[i].ID == ID {
			return &s.Devices[i]
		}
	}
	return nil
}

func (s *Site) GetDeviceByServiceAddress(addr string)*Device {
	for i := range s.Devices {
		for _,v := range s.Devices[i].Service {
			if v.Addr == addr {
				return &s.Devices[i]
			}
		}
	}
	return nil
}

func (s *Site) GetServiceByAddress(addr string) *Service {
	for i := range s.Devices {
		for _,v := range s.Devices[i].Service {
			if v.Addr == addr {
				return &v
			}
		}
	}
	return nil
}