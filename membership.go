package asana

// Membership stores Membership from Asana
//
type Membership struct {
	Project CompactObject `json:"project"`
	Section CompactObject `json:"section"`
}
