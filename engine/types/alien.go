package types


import (
	"fmt"
)

// Alien Type definition
type Alien struct {
	// Alien ID
	AlienID int
	// City where Alien resides in
	City *City
	//Flag whether Alien is trapped
	IsTrapped bool
}

// Generate New Alien
func NewAlien(alienID int) *Alien {
	return &Alien{
		AlienID: alienID,
	}
}

// String output of Alien
func (a *Alien) String() string {
	return fmt.Sprintf("Alien #%d", a.AlienID)
}
