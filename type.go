package main

type AlmanaxEvent struct {
	ItemImage string `json:"itemImage"`
	Quest     string `json:"quest"`
	Type      string `json:"type"`
	Effect    string `json:"effect"`
	Offering  string `json:"offering"`
}

// AlmanaxCalendar the map key is the date (Format MM/DD)
type AlmanaxCalendar map[string]AlmanaxEvent
