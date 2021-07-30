package asana

// Membership stores Membership from Asana
//
type Membership struct {
	Project Object `json:"project"`
	Section Object `json:"section"`
}
