package inventory

import "encoding/json"

// Color holds primary colors.
type Color struct {
	R int
	G int
	B int
}

// Tab describes a tabulation style (background color,
// custom name, style, and so on)
type Tab struct {
	Name            string `json:"n"`
	Index           int    `json:"i"`
	Id              string `json:"id"`
	Type            string `json:"type"`
	Hidden          bool   `json:"hidden"`
	Selected        bool   `json:"selected"`
	BackgroundColor Color  `json:"colour"`
	ImgL            string `json:"srcL"`
	ImgC            string `json:"srcC"`
	ImgR            string `json:"srcR"`
}

// StashTab holds all stash tabulations (thus all items).
type StashTab struct {
	NumTabs int
	Tabs    []Tab
	Items   []Item
	// CurrencyLayout FIXME
}

func (s *StashTab) String() string {
	json, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return "<marshalling error>"
	}
	return string(json)
}
