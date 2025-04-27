package models

type Person struct {
	Id    int    `json:"-"`
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	//Memberships []Subscription `json:"memberships,omitempty" required:"false"`
}
