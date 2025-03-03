package models

type Person struct {
	Id   int    `json:"-"`
	Name string `json:"name,omitempty"`
	//Memberships []Membership `json:"memberships,omitempty" required:"false"`
}
