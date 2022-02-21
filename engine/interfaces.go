package engine

import (
	"context"

	"alien-invasion-cc/engine/types"
)

// WorldStorer is a world store interface
type World interface {
	// GetCity retrieves a city
	GetCity(ctx context.Context, cityName string) (*types.City, error)
	// GetAliveCities retrieves the list of non destroyed cities
	GetAliveCities(ctx context.Context) ([]*types.City, error)
	// AddCity adds a city
	AddCity(ctx context.Context, cityName string) (*types.City, error)
	// DestroyCity destroys a city
	DestroyCity(ctx context.Context, city *types.City) error
	// AddLink adds a link from a city to another city given a direction
	AddLink(ctx context.Context, cityFrom, cityTo *types.City, direction types.Direction) error
	// GetAlien retrieves an alien
	GetAlien(ctx context.Context, alienID int) (*types.Alien, error)
	// AddAlien adds an alien
	AddAlien(ctx context.Context, alienID int) (*types.Alien, error)
	// MoveAlien moves an alien to a city
	MoveAlien(ctx context.Context, alien *types.Alien, city *types.City) error
	// IsTrappedAlien checks if an alien is trapped
	IsTrappedAlien(ctx context.Context, alien *types.Alien) (bool, error)
	// TrapAlien traps an alien
	TrapAlien(ctx context.Context, alien *types.Alien) error
	// GetAlienAtCity retrieves the alien at a given city
	GetAlienAtCity(ctx context.Context, city *types.City) (*types.Alien, error)
	// GetUntrappedAliens retrieves the list of untrapped aliens
	GetUntrappedAliens(ctx context.Context) ([]*types.Alien, error)
}

// Simulator is an alien invasion simulator interface
type Engine interface {
	// LoadEngine prepares the simulation
	LoadEngine(ctx context.Context) error
	// HasNextMove computes if a next step of the simulation exists
	HasNextMove(ctx context.Context) (bool, error)
	// DoNextMove simulates the next step of the simulation
	DoNextMove(ctx context.Context) error
	// Run simulates an alien invasion
	Run(ctx context.Context) error
	// Finalize finalizes the simulation
	Finalize(ctx context.Context) error
}