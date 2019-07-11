package primefimp


type Site struct {
	Id        int        `json:"id,omitempty"`
	Devices   []Device   `json:"device,omitempty"`
	Things    []Thing    `json:"thing,omitempty"`
	Rooms     []Room     `json:"room,omitempty"`
	House     *House     `json:"house,omitempty"`
	Hub       *Hub       `json:"hub,omitempty"`
	Areas     []Area     `json:"area,omitempty"`
	Shortcuts []Shortcut `json:"shortcut,omitempty"`
	Problem   bool       `json:"problem,omitempty"`
}

func SiteFromResponse(resp *Response) *Site {
	site := Site{Devices: resp.GetDevices(), Things:resp.GetThings(), Rooms:resp.GetRooms(), Areas:resp.GetAreas(),House:resp.GetHouse(),
		Shortcuts: resp.GetShortcuts()}
	return &site

}
