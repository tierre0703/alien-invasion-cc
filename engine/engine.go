package engine

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"alien-invasion-cc/engine/types"
)

// Alien Invasion Engine type definition
type EngineImpl struct {
	world World

	in io.Reader

	out io.Writer
	
	numAliens uint

	maxMoves uint

	totalMoves uint

}

var _ Engine = (*EngineImpl)(nil)

// Generate Random Int
func GetRandInt(n int) (int, error) {
	rand.Seed(time.Now().UnixNano())
	r := 0
	if n <= 0 {
		return r, types.ERR_RANDOM_OUT_OF_BOUNDS
	}

	r = rand.Intn(n)
	return r, nil
}

func NewEngine(numAliens, maxMoves uint, in io.Reader, out io.Writer) *EngineImpl {
	world := NewWorld()
	return &EngineImpl{
		world: 		world,
		in: 		in,
		out:		out,
		maxMoves:	maxMoves,
		numAliens:numAliens,
	}
}

// LoadEngine - spawn aliens, load world
func (s *EngineImpl) LoadEngine(ctx context.Context) error {

	err := s.loadWorld(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < int(s.numAliens); i++ {
		alienID := i + 1
		alien, err := s.world.AddAlien(ctx, alienID)
		if err != nil {
			return err
		}

		var nextCity *types.City
		aliveCities, err := s.world.GetAliveCities(ctx)
		if err != nil {
			return err
		}

		if len(aliveCities) == 0 {
			return nil
		}

		r, err := GetRandInt(len(aliveCities))
		if err != nil {
			return err
		}

		nextCity = aliveCities[r]
		_, err = s.moveAlienToCity(ctx, alien, nextCity)
		if err != nil {
			return err
		}

	}

	return nil
}

// HasNextMove check if next move available 
func (s *EngineImpl) HasNextMove(ctx context.Context) (bool, error) {

	if s.totalMoves >= s.maxMoves {
		return false, nil
	}

	untrappedAliens, err := s.world.GetUntrappedAliens(ctx)
	if err != nil {
		return false, err
	}

	if len(untrappedAliens) == 0 {
		return false, nil
	}

	aliveCities, err := s.world.GetAliveCities(ctx)
	if err != nil {
		return false, err
	}

	if len(aliveCities) == 0 {
		return false, nil
	}

	return true, nil
}

// DoNextMove proceed next move of engine
func (s *EngineImpl) DoNextMove(ctx context.Context) error {

	s.totalMoves++
	untrappedAliens, err :=s.world.GetUntrappedAliens(ctx)
	if err != nil {
		return err
	}

	for _, alien := range untrappedAliens {

		isTrapped, err := s.world.IsTrappedAlien(ctx, alien)
		if err != nil {
			return err
		}

		if isTrapped {
			continue
		}

		var nextCity *types.City
		currentCity := alien.City
		availableLinks := currentCity.GetAvailableLinks()
		if len(availableLinks) > 0 {
			r, err := GetRandInt(len(availableLinks))
			if err != nil {
				return err
			}

			i := 0
			for _, city := range availableLinks {
				if i == r {
					nextCity = city
					continue
				}
				i++
			}

			_, err = s.moveAlienToCity(ctx, alien, nextCity)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func run(ctx context.Context, s Engine) error {

	err := s.LoadEngine(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Warn("Simulation cancelled")
			return types.ERR_CONTEXT_CANCELLED
		
		default:
			hasNextStep, err := s.HasNextMove(ctx)
			if err != nil {
				return err
			}

			if !hasNextStep {
				return s.Finalize(ctx)
			}

			err = s.DoNextMove(ctx)
			if err != nil {
				return err
			}
		}
	}
}


func (s *EngineImpl) Run (ctx context.Context) error {
	return run(ctx, s)
}

// Finalize engine finalize and output result
func (s *EngineImpl) Finalize(ctx context.Context) error {

	fmt.Fprintf(s.out, "\n===================\n")
	fmt.Fprintf(s.out, "Simulation Finished\n")
	fmt.Fprintf(s.out, "===================\n")

	cities, err := s.world.GetAliveCities(ctx)

	if err != nil {
		return err
	}

	fmt.Fprintf(s.out, "Remain Cities: %d\n", len(cities))

	_, err = fmt.Fprintln(s.out, "")
	if err != nil {
		return err
	}

	for _, city := range cities {
		_, err = fmt.Fprintln(s.out, city)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadWorld load City and City Link from Map Data
func (s * EngineImpl) loadWorld(ctx context.Context) error {

	registerCity := func(cityName string)(*types.City, error) {
		city, err := s.world.GetCity(ctx, cityName)
		if err != nil {
			return city, err
		}

		if city == nil {
			city, err = s.world.AddCity(ctx, cityName)
			if err != nil {
				return city, err
			}
		}
		return city, nil
	}

	scanner := bufio.NewScanner(s.in)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		lineChunks := strings.Split(line, " ")
		if len(lineChunks) == 0 {
			return types.ERR_PARSE_CITY_DEFINITION
		}

		cityFromName := lineChunks[0]
		cityFrom, err := registerCity(cityFromName)
		if err != nil {
			return err
		}

		for _, lineChunk := range lineChunks[1:] {
			linkChunks := strings.Split(strings.TrimSpace(lineChunk), "=")
			if len(linkChunks) != 2 {
				return types.ERR_PARSE_CITY_DEFINITION
			}

			directionName := linkChunks[0]
			cityToName := linkChunks[1]
			cityTo, err := registerCity(cityToName)
			if err != nil {
				return err
			}

			var direction types.Direction
			switch directionName {
			case "north":
				direction = types.North
			case "east":
				direction = types.East
			case "south":
				direction = types.South
			case "west":
				direction = types.West
			default:
				return types.ERR_PARSE_CITY_DEFINITION
			}
			err = s.world.AddLink(ctx, cityFrom, cityTo, direction)
			if err != nil {
				return err
			}
		
		}
	}

	return nil
}

//moveAlienToCity move Alien to target city
func (s *EngineImpl) moveAlienToCity(ctx context.Context, alien *types.Alien, city *types.City) (bool, error) {

	destroyedCity := false
	alienAlreadyInCity, err := s.world.GetAlienAtCity(ctx, city)
	if err != nil {
		return destroyedCity, err
	}

	switch {
	case alienAlreadyInCity == alien:
		return destroyedCity, nil
	case alienAlreadyInCity == nil:
		err := s.world.MoveAlien(ctx, alien, city)
		if err != nil {
			return destroyedCity, err
		}
	default:
		err = s.world.TrapAlien(ctx, alien)
		if err != nil {
			return destroyedCity, err
		}
		
		err = s.world.TrapAlien(ctx, alienAlreadyInCity)
		if err != nil {
			return destroyedCity, err
		}

		err = s.world.DestroyCity(ctx, city)
		if err != nil {
			return destroyedCity, err
		}

		destroyedCity = true
		_, err := fmt.Fprintf(s.out, "%s has been destroyed by %s and %s\n", city.Name, alien, alienAlreadyInCity)
		if err != nil {
			return destroyedCity, err
		}
	}

	return destroyedCity, nil
}