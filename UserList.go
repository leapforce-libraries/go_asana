package asana

// UserList stores UserList from Asana
//
type UserList struct {
	ID   string        `json:"gid"`
	User CompactObject `json:"user"`
}
