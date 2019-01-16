package cerebrum

// Company represents company model
type Company struct {
	Base
	Name      string     `json:"name"`
	Locations []Location `json:"locations,omitempty"`
	Owner     User       `json:"owner"`
}
