package types

import (
	"fmt"
	"strings"
)

// City Type Definition

type City struct {
	Name string
	North, East, South, West *City
}

// City constructor
func NewCity(Name string) *City {
	return &City{
		Name: Name,
	}
}

// GetCityLink retrieves the destination city given a direction
func (c *City) GetCityLink(direction Direction) (*City, error) {

	switch direction {
	case North:
		return c.North, nil
	case East:
		return c.East, nil
	case South:
		return c.South, nil
	case West:
		return c.West, nil
	default:
		var city *City
		return city, ERR_UNKNOWN_DIRECTION
	}
}


// SetCityLink sets destination city to given direction
func (c *City) SetCityLink(city *City, direction Direction) error {
	switch direction {
	case North:
		c.North = city
	case East:
		c.East = city
	case South:
		c.South = city
	case West:
		c.West = city
	default:
		return ERR_UNKNOWN_DIRECTION
	}

	return nil
}


// RemoveCityLink removes the destination city in any direction if it exists
func (c *City) RemoveCityLink(city *City) error {
	switch {
	case c.North == city:
		c.North = nil
	case c.South == city:
		c.South = nil
	case c.East == city:
		c.East = nil
	case c.West == city:
		c.West = nil
	default:
		return ERR_UNKNOWN_CITY
	}

	return nil
}

// GetAvailableLinks retrieves the available links from this city
func (c *City) GetAvailableLinks() map[Direction]*City {
	links := make(map[Direction]*City)
	if c.North != nil {
		links[North] = c.North
	}
	if c.South != nil {
		links[South] = c.South
	}
	if c.West != nil {
		links[West] = c.West
	}
	if c.East != nil {
		links[East] = c.East
	}
	return links
}

// String output of City
func (c *City) String() string {
	chunks := []string{c.Name}
	if c.North != nil {
		chunks = append(chunks, fmt.Sprintf("north=%s", c.North.Name))
	}
	if c.East != nil {
		chunks = append(chunks, fmt.Sprintf("east=%s", c.East.Name))
	}
	if c.South != nil {
		chunks = append(chunks, fmt.Sprintf("south=%s", c.South.Name))
	}
	if c.West != nil {
		chunks = append(chunks, fmt.Sprintf("west=%s", c.West.Name))
	}
	return strings.Join(chunks, " ")
}