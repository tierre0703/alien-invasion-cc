package engine

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"alien-invasion-cc/engine/types"
	"github.com/stretchr/testify/require"
)


func Test_Engine_HasNextMove(t *testing.T) {
	alien1 := types.NewAlien(1)
	alien2 := types.NewAlien(2)
	city1 := types.NewCity("City1")
	city2 := types.NewCity("City2")

	aliensEmpty := []*types.Alien{}
	aliensFilled := []*types.Alien{alien1, alien2}
	citiesEmpty := []*types.City{}
	citiesFilled := []*types.City{city1, city2}

	error1 := fmt.Errorf("error 1")
	error2 := fmt.Errorf("error 2")

	tests := []struct {
		testName                     string
		giveTotalSteps, giveMaxSteps uint
		giveUntrappedAliens          []*types.Alien
		giveUntrappedAliensError     error
		giveAliveCities              []*types.City
		giveAliveCitiesError         error
		wantGetUntrappedAliensCalls  int
		wantGetAliveCitiesCalls      int
		wantResult                   bool
		wantError                    error
	}{
		{
			testName:                    "Too many steps",
			giveTotalSteps:              10,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensEmpty,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesEmpty,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 0,
			wantGetAliveCitiesCalls:     0,
			wantResult:                  false,
			wantError:                   nil,
		},
		{
			testName:                    "All aliens trapped",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensEmpty,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesFilled,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     0,
			wantResult:                  false,
			wantError:                   nil,
		},
		{
			testName:                    "GetUntrappedAliens returns error",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensEmpty,
			giveUntrappedAliensError:    error1,
			giveAliveCities:             citiesFilled,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     0,
			wantResult:                  false,
			wantError:                   error1,
		},
		{
			testName:                    "All cities destroyed",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensFilled,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesEmpty,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     1,
			wantResult:                  false,
			wantError:                   nil,
		},
		{
			testName:                    "GetAliveCities returns error",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensFilled,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesEmpty,
			giveAliveCitiesError:        error2,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     1,
			wantResult:                  false,
			wantError:                   error2,
		},
		{
			testName:                    "Next step exists",
			giveTotalSteps:              2,
			giveMaxSteps:                5,
			giveUntrappedAliens:         aliensFilled,
			giveUntrappedAliensError:    nil,
			giveAliveCities:             citiesFilled,
			giveAliveCitiesError:        nil,
			wantGetUntrappedAliensCalls: 1,
			wantGetAliveCitiesCalls:     1,
			wantResult:                  true,
			wantError:                   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctx := context.Background()

			worldMock := &WorldMock{}
			if tt.wantGetUntrappedAliensCalls > 0 {
				worldMock.On("GetUntrappedAliens", ctx).Return(tt.giveUntrappedAliens, tt.giveUntrappedAliensError).Times(tt.wantGetUntrappedAliensCalls)

			}
			if tt.wantGetAliveCitiesCalls > 0 {
				worldMock.On("GetAliveCities", ctx).Return(tt.giveAliveCities, tt.giveAliveCitiesError).Times(tt.wantGetAliveCitiesCalls)
			}
			defer worldMock.AssertExpectations(t)

			s := EngineImpl{
				world:       worldMock,
				in:          &bytes.Buffer{},
				out:         &bytes.Buffer{},
				totalMoves:  tt.giveTotalSteps,
				maxMoves:    tt.giveMaxSteps,
				numAliens: 0,
			}

			result, err := s.HasNextMove(ctx)
			require.ErrorIs(t, tt.wantError, err)
			require.Equal(t, tt.wantResult, result)
		})
	}
}

func Test_Engine_DoNextMove(t *testing.T) {
	var alienNil *types.Alien
	alien1 := types.NewAlien(1)
	alien2 := types.NewAlien(2)
	alien3 := types.NewAlien(3)
	city1 := types.NewCity("City1")
	city2 := types.NewCity("City2")

	t.Run("Case 1: Alien1 is trapped / Alien2 is moving to same city", func(t *testing.T) {
		ctx := context.Background()

		err := city1.SetCityLink(city2, types.South)
		require.NoError(t, err)
		alien2.City = city1

		worldMock := &WorldMock{}
		worldMock.On("GetUntrappedAliens", ctx).Return([]*types.Alien{alien1, alien2}, nil).Once()
		// Alien1 is already trapped
		worldMock.On("IsTrappedAlien", ctx, alien1).Return(true, nil).Once()
		// Alien2 is moved to its current city
		worldMock.On("IsTrappedAlien", ctx, alien2).Return(false, nil).Once()
		worldMock.On("GetAlienAtCity", ctx, city2).Return(alien2, nil).Once()
		defer worldMock.AssertExpectations(t)

		s := EngineImpl{
			world:       worldMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err = s.DoNextMove(ctx)
		require.NoError(t, err)
		require.Equal(t, uint(1), s.totalMoves)
	})

	t.Run("Case 2: Alien1 is trapped / Alien2 is moving to an unoccupied city", func(t *testing.T) {
		ctx := context.Background()

		err := city1.SetCityLink(city2, types.South)
		require.NoError(t, err)
		alien2.City = city1

		worldMock := &WorldMock{}
		worldMock.On("GetUntrappedAliens", ctx).Return([]*types.Alien{alien1, alien2}, nil).Once()
		// Alien1 is already trapped
		worldMock.On("IsTrappedAlien", ctx, alien1).Return(true, nil).Once()
		// Alien2 is moved to an unoccupied city
		worldMock.On("IsTrappedAlien", ctx, alien2).Return(false, nil).Once()
		worldMock.On("GetAlienAtCity", ctx, city2).Return(alienNil, nil).Once()
		worldMock.On("MoveAlien", ctx, alien2, city2).Return(nil).Once()
		defer worldMock.AssertExpectations(t)

		s := EngineImpl{
			world:       worldMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err = s.DoNextMove(ctx)
		require.NoError(t, err)
		require.Equal(t, uint(1), s.totalMoves)
	})

	t.Run("Case 3:  Alien1 is trapped / Alien2 is moving to an occupied city", func(t *testing.T) {
		ctx := context.Background()

		err := city1.SetCityLink(city2, types.South)
		require.NoError(t, err)
		alien2.City = city1

		worldMock := &WorldMock{}
		worldMock.On("GetUntrappedAliens", ctx).Return([]*types.Alien{alien1, alien2}, nil).Once()
		// Alien1 is already trapped
		worldMock.On("IsTrappedAlien", ctx, alien1).Return(true, nil).Once()
		// Alien2 is moved to an occupied city
		worldMock.On("IsTrappedAlien", ctx, alien2).Return(false, nil).Once()
		worldMock.On("GetAlienAtCity", ctx, city2).Return(alien3, nil).Once()
		worldMock.On("TrapAlien", ctx, alien2).Return(nil).Once()
		worldMock.On("TrapAlien", ctx, alien3).Return(nil).Once()
		worldMock.On("DestroyCity", ctx, city2).Return(nil).Once()
		defer worldMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := EngineImpl{
			world:       worldMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err = s.DoNextMove(ctx)
		require.NoError(t, err)
		require.Equal(t, uint(1), s.totalMoves)
		require.Equal(t, "City2 has been destroyed by Alien #2 and Alien #3\n", out.String())
	})

	t.Run("Case 2: Error", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldMock := &WorldMock{}
		worldMock.On("GetUntrappedAliens", ctx).Return([]*types.Alien{nil}, error1).Once()
		defer worldMock.AssertExpectations(t)

		s := EngineImpl{
			world:       worldMock,
			in:          &bytes.Buffer{},
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.DoNextMove(ctx)
		require.ErrorIs(t, err, error1)
	})
}

func Test_Engine_Run(t *testing.T) {
	totalSteps := 10
	error1 := fmt.Errorf("error 1")

	t.Run("Case 1: OK", func(t *testing.T) {
		ctx := context.Background()

		engineMock := &EngineMock{}
		engineMock.On("LoadEngine", ctx).Return(nil).Once()
		engineMock.On("HasNextMove", ctx).Return(true, nil).Times(totalSteps)
		engineMock.On("HasNextMove", ctx).Return(false, nil).Once()
		engineMock.On("DoNextMove", ctx).Return(nil).Times(totalSteps)
		engineMock.On("Finalize", ctx).Return(nil).Once()
		defer engineMock.AssertExpectations(t)

		err := run(ctx, engineMock)
		require.NoError(t, err)
	})

	t.Run("Case 2: Error on Prepare", func(t *testing.T) {
		ctx := context.Background()

		engineMock := &EngineMock{}
		engineMock.On("LoadEngine", ctx).Return(error1).Once()
		defer engineMock.AssertExpectations(t)

		err := run(ctx, engineMock)
		require.ErrorIs(t, err, error1)
	})

	t.Run("Case 3: Error on HasNextMove", func(t *testing.T) {
		ctx := context.Background()

		engineMock := &EngineMock{}
		engineMock.On("LoadEngine", ctx).Return(nil).Once()
		engineMock.On("HasNextMove", ctx).Return(true, nil).Times(totalSteps)
		engineMock.On("HasNextMove", ctx).Return(false, error1).Once()
		engineMock.On("DoNextMove", ctx).Return(nil).Times(totalSteps)
		defer engineMock.AssertExpectations(t)

		err := run(ctx, engineMock)
		require.ErrorIs(t, err, error1)
	})

	t.Run("Case 4: Error on DoNextMove", func(t *testing.T) {
		ctx := context.Background()

		engineMock := &EngineMock{}
		engineMock.On("LoadEngine", ctx).Return(nil).Once()
		engineMock.On("HasNextMove", ctx).Return(true, nil).Times(totalSteps)
		engineMock.On("DoNextMove", ctx).Return(nil).Times(totalSteps - 1)
		engineMock.On("DoNextMove", ctx).Return(error1).Once()
		defer engineMock.AssertExpectations(t)

		err := run(ctx, engineMock)
		require.ErrorIs(t, err, error1)
	})

	t.Run("Case 5: Error on Finalize", func(t *testing.T) {
		ctx := context.Background()

		engineMock := &EngineMock{}
		engineMock.On("LoadEngine", ctx).Return(nil).Once()
		engineMock.On("HasNextMove", ctx).Return(true, nil).Times(totalSteps)
		engineMock.On("HasNextMove", ctx).Return(false, nil).Once()
		engineMock.On("DoNextMove", ctx).Return(nil).Times(totalSteps)
		engineMock.On("Finalize", ctx).Return(error1).Once()
		defer engineMock.AssertExpectations(t)

		err := run(ctx, engineMock)
		require.ErrorIs(t, err, error1)
	})
}

func Test_Engine_Finalize(t *testing.T) {
	city1 := types.NewCity("City1")
	city2 := types.NewCity("City2")

	t.Run("Case 1: OK", func(t *testing.T) {
		ctx := context.Background()

		worldMock := &WorldMock{}
		worldMock.On("GetAliveCities", ctx).Return([]*types.City{city1, city2}, nil).Once()
		defer worldMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := EngineImpl{
			world:       worldMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.Finalize(ctx)
		require.NoError(t, err)
		require.Equal(t, "\nCity1\nCity2\n", out.String())
	})

	t.Run("Case 2: Error", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldMock := &WorldMock{}
		worldMock.On("GetAliveCities", ctx).Return([]*types.City{nil}, error1).Once()
		defer worldMock.AssertExpectations(t)

		out := &bytes.Buffer{}

		s := EngineImpl{
			world:       worldMock,
			in:          &bytes.Buffer{},
			out:         out,
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.Finalize(ctx)
		require.ErrorIs(t, err, error1)
		require.Equal(t, "", out.String())
	})
}

func Test_Engine_loadWorld(t *testing.T) {
	var cityNil *types.City
	city1 := types.NewCity("City1")
	city2 := types.NewCity("City2")
	city3 := types.NewCity("City3")
	city4 := types.NewCity("City4")
	city5 := types.NewCity("City5")
	city6 := types.NewCity("City6")
	city7 := types.NewCity("City7")

	t.Run("Case 1: Empty file", func(t *testing.T) {
		ctx := context.Background()

		worldMock := &WorldMock{}
		defer worldMock.AssertExpectations(t)

		input := ""
		in := strings.NewReader(input)

		s := EngineImpl{
			world:       worldMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.NoError(t, err)
	})

	t.Run("Case 1: Correctly formed file", func(t *testing.T) {
		ctx := context.Background()

		worldMock := &WorldMock{}
		// City1
		worldMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldMock.On("GetCity", ctx, "City1").Return(city1, nil).Once()
		worldMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		// City2
		worldMock.On("GetCity", ctx, "City2").Return(cityNil, nil).Once()
		worldMock.On("GetCity", ctx, "City2").Return(city2, nil).Once()
		worldMock.On("AddCity", ctx, "City2").Return(city2, nil).Once()
		// City3
		worldMock.On("GetCity", ctx, "City3").Return(cityNil, nil).Once()
		worldMock.On("GetCity", ctx, "City3").Return(city3, nil).Times(2)
		worldMock.On("AddCity", ctx, "City3").Return(city3, nil).Once()
		// City4
		worldMock.On("GetCity", ctx, "City4").Return(cityNil, nil).Once()
		worldMock.On("GetCity", ctx, "City4").Return(city4, nil).Times(2)
		worldMock.On("AddCity", ctx, "City4").Return(city4, nil).Once()
		// City5
		worldMock.On("GetCity", ctx, "City5").Return(cityNil, nil).Once()
		worldMock.On("GetCity", ctx, "City5").Return(city5, nil).Times(3)
		worldMock.On("AddCity", ctx, "City5").Return(city5, nil).Once()
		// City6
		worldMock.On("GetCity", ctx, "City6").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City6").Return(city6, nil).Once()
		// City7
		worldMock.On("GetCity", ctx, "City7").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City7").Return(city7, nil).Once()
		// Links from City1
		worldMock.On("AddLink", ctx, city1, city2, types.North).Return(nil).Once()
		worldMock.On("AddLink", ctx, city1, city3, types.East).Return(nil).Once()
		worldMock.On("AddLink", ctx, city1, city4, types.South).Return(nil).Once()
		worldMock.On("AddLink", ctx, city1, city5, types.West).Return(nil).Once()
		// Links from City2
		worldMock.On("AddLink", ctx, city2, city1, types.East).Return(nil).Once()
		worldMock.On("AddLink", ctx, city2, city4, types.South).Return(nil).Once()
		// Links from City3
		worldMock.On("AddLink", ctx, city3, city5, types.West).Return(nil).Once()
		worldMock.On("AddLink", ctx, city3, city7, types.East).Return(nil).Once()
		worldMock.On("AddLink", ctx, city3, city5, types.South).Return(nil).Once()
		// Links from City4
		worldMock.On("AddLink", ctx, city4, city3, types.East).Return(nil).Once()
		worldMock.On("AddLink", ctx, city4, city5, types.North).Return(nil).Once()
		defer worldMock.AssertExpectations(t)

		input := `
City1 north=City2 east=City3 south=City4 west=City5
City2 east=City1 south=City4
	City3 west=City5 east=City7 south=City5
City4 east=City3 north=City5
City6

		`
		in := strings.NewReader(input)

		s := EngineImpl{
			world:       worldMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.NoError(t, err)
	})

	t.Run("Case 3: Incorrect direction", func(t *testing.T) {
		ctx := context.Background()

		worldMock := &WorldMock{}
		// City1
		worldMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		// City2
		worldMock.On("GetCity", ctx, "City2").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City2").Return(city2, nil).Once()
		defer worldMock.AssertExpectations(t)

		input := `
City1 test=City2	
City2 east=City1 south=City4
		`
		in := strings.NewReader(input)

		s := EngineImpl{
			world:       worldMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.ErrorIs(t, err, types.ERR_PARSE_CITY_DEFINITION)
	})

	t.Run("Case 3: Incorrect format", func(t *testing.T) {
		ctx := context.Background()

		worldMock := &WorldMock{}
		// City1
		worldMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		defer worldMock.AssertExpectations(t)

		input := `
City1 test=City2=City3	
City2 east=City1 south=City4
		`
		in := strings.NewReader(input)

		s := EngineImpl{
			world:       worldMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.ErrorIs(t, err, types.ERR_PARSE_CITY_DEFINITION)
	})

	t.Run("Case 4: Error in AddCity", func(t *testing.T) {
		ctx := context.Background()

		error1 := fmt.Errorf("error 1")

		worldMock := &WorldMock{}
		// City1
		worldMock.On("GetCity", ctx, "City1").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City1").Return(city1, nil).Once()
		// City2
		worldMock.On("GetCity", ctx, "City2").Return(cityNil, nil).Once()
		worldMock.On("AddCity", ctx, "City2").Return(city2, error1).Once()
		defer worldMock.AssertExpectations(t)

		input := `
City1 north=City2	
City2 east=City1 south=City4
		`
		in := strings.NewReader(input)

		s := EngineImpl{
			world:       worldMock,
			in:          in,
			out:         &bytes.Buffer{},
			totalMoves:  0,
			maxMoves:    10,
			numAliens: 0,
		}

		err := s.loadInputToWorld(ctx)
		require.ErrorIs(t, err, error1)
	})
}