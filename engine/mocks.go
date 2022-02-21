package engine

import (
	"context"

	"github.com/stretchr/testify/mock"

	"alien-invasion-cc/engine/types"
)

// WorldMock mocks a WorldStorer
type WorldMock struct {
	mock.Mock
}

var _ World = (*WorldMock)(nil)

// GetCity retrieves a city
func (w *WorldMock) GetCity(ctx context.Context, cityName string) (*types.City, error) {
	args := w.Called(ctx, cityName)
	return args.Get(0).(*types.City), args.Error(1)
}

// GetAliveCities retrieves the list of non destroyed cities
func (w *WorldMock) GetAliveCities(ctx context.Context) ([]*types.City, error) {
	args := w.Called(ctx)
	return args.Get(0).([]*types.City), args.Error(1)
}

// AddCity adds a city
func (w *WorldMock) AddCity(ctx context.Context, cityName string) (*types.City, error) {
	args := w.Called(ctx, cityName)
	return args.Get(0).(*types.City), args.Error(1)
}

// DestroyCity destroys a city
func (w *WorldMock) DestroyCity(ctx context.Context, city *types.City) error {
	args := w.Called(ctx, city)
	return args.Error(0)
}

// AddLink adds a link from a city to another city given a direction
func (w *WorldMock) AddLink(ctx context.Context, cityFrom, cityTo *types.City, direction types.Direction) error {
	args := w.Called(ctx, cityFrom, cityTo, direction)
	return args.Error(0)
}

// GetAlien retrieves an alien
func (w *WorldMock) GetAlien(ctx context.Context, alienID int) (*types.Alien, error) {
	args := w.Called(ctx, alienID)
	return args.Get(0).(*types.Alien), args.Error(1)
}

// AddAlien adds an alien
func (w *WorldMock) AddAlien(ctx context.Context, alienID int) (*types.Alien, error) {
	args := w.Called(ctx, alienID)
	return args.Get(0).(*types.Alien), args.Error(1)
}

// MoveAlien moves an alien to a city
func (w *WorldMock) MoveAlien(ctx context.Context, alien *types.Alien, city *types.City) error {
	args := w.Called(ctx, alien, city)
	return args.Error(0)
}

// IsTrappedAlien checks if an alien is trapped
func (w *WorldMock) IsTrappedAlien(ctx context.Context, alien *types.Alien) (bool, error) {
	args := w.Called(ctx, alien)
	return args.Bool(0), args.Error(1)
}

// TrapAlien traps an alien
func (w *WorldMock) TrapAlien(ctx context.Context, alien *types.Alien) error {
	args := w.Called(ctx, alien)
	return args.Error(0)
}

// GetAlienAtCity retrieves the alien at a given city
func (w *WorldMock) GetAlienAtCity(ctx context.Context, city *types.City) (*types.Alien, error) {
	args := w.Called(ctx, city)
	return args.Get(0).(*types.Alien), args.Error(1)
}

// GetUntrappedAliens retrieves the list of untrapped aliens
func (w *WorldMock) GetUntrappedAliens(ctx context.Context) ([]*types.Alien, error) {
	args := w.Called(ctx)
	return args.Get(0).([]*types.Alien), args.Error(1)
}

// EngineMock mocks a Simulator
type EngineMock struct {
	mock.Mock
}

var _ Engine = (*EngineMock)(nil)

// Prepare prepares the simulation
func (s *EngineMock) LoadEngine(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// HasNextMove computes if a next step of the simulation exists
func (s *EngineMock) HasNextMove(ctx context.Context) (bool, error) {
	args := s.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// DoNextMove simulates the next step of the simulation
func (s *EngineMock) DoNextMove(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// Run simulates an alien invasion
func (s *EngineMock) Run(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}

// Finalize finalizes the simulation
func (s *EngineMock) Finalize(ctx context.Context) error {
	args := s.Called(ctx)
	return args.Error(0)
}
