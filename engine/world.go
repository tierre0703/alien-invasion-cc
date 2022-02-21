package engine


import (
	"context"
	"alien-invasion-cc/engine/types"
)


type WorldImpl struct {

	cities map[string]*types.City

	aliens map[int]*types.Alien

	alienInCities map[*types.City]*types.Alien

	links map[*types.City][]*types.City
}

var _ World = (*WorldImpl)(nil)

// World Constructor
func NewWorld() *WorldImpl {
	var (
		cities	 			= make(map[string]*types.City)
		aliens				= make(map[int]*types.Alien)
		alienInCities		= make(map[*types.City]*types.Alien)
		links	= make(map[*types.City][]*types.City) 
	)

	return &WorldImpl{
		cities: 			cities,
		aliens:				aliens,
		alienInCities: 		alienInCities,
		links:	links,
	}
}

// GetCity retrieves city with name
func (w *WorldImpl) GetCity(ctx context.Context, cityName string) (*types.City, error) {

	var city *types.City
	if cityFound,found := w.cities[cityName]; found {
		return cityFound, nil
	}

	return city, nil
}

// AddCity add a city to the world
func (w *WorldImpl) AddCity(ctx context.Context, cityName string) (*types.City, error) {

	var city *types.City
	if cityName == "" {
		return city, types.ERR_EMPTY_CITY_NAME
	}

	if _, found := w.cities[cityName]; found {
		return city, types.ERR_DUPLICATE_CITY
	}

	newCity := types.NewCity(cityName)
	w.cities[newCity.Name] = newCity

	return newCity, nil
}

// DestroyCity remove city from world
func (w *WorldImpl) DestroyCity(ctx context.Context, city *types.City) error {

	if citiesFrom, found := w.links[city]; found {
		for _, cityFrom := range citiesFrom {
			err := cityFrom.RemoveCityLink(city)
			if err != nil {
				return err
			}
		}
	}

	delete(w.cities, city.Name)
	delete(w.alienInCities, city)
	delete(w.links, city)
	
	return nil
}

// GetAliveCities retrieves list of non-destroyed cities
func (w *WorldImpl) GetAliveCities(ctx context.Context) ([]*types.City, error) {

	var cities []*types.City
	for _, city := range w.cities {
		cities = append(cities, city)
	}

	return cities, nil
}

// AddLink add a link from a city to another city with direction
func (w *WorldImpl) AddLink(ctx context.Context, cityFrom, cityTo *types.City, direction types.Direction) error {

	if cityFrom == nil {
		return types.ERR_MISSING_CITY
	}

	if cityTo == nil {
		return types.ERR_MISSING_CITY
	}

	if cityFrom.Name == cityTo.Name {
		return types.ERR_LINK_SAME_CITY
	}

	cityFromFound, err := w.GetCity(ctx, cityFrom.Name)
	if err != nil {
		return err
	}

	if cityFromFound == nil {
		return types.ERR_UNKNOWN_CITY
	}

	cityToFound, err := w.GetCity(ctx, cityTo.Name)
	if err != nil {
		return err
	}

	if cityToFound == nil {
		return types.ERR_UNKNOWN_CITY
	}

	cityToRegistered, err := cityFrom.GetCityLink(direction)
	if err != nil {
		return err
	}

	if cityToRegistered != nil && cityToRegistered != cityTo {
		return types.ERR_ALREADY_EXISTS_LINK
	}

	err = cityFrom.SetCityLink(cityTo, direction)
	if err != nil {
		return err
	}

	citiesFrom, found := w.links[cityTo]
	if !found {
		citiesFrom = make([]*types.City, 0)
	}

	citiesFrom = append(citiesFrom, cityFrom)
	w.links[cityTo] = citiesFrom
	return nil
}

// GetAlien retrieves Alien by alienID
func (w *WorldImpl) GetAlien(ctx context.Context, alienID int) (*types.Alien, error) {

	var alien *types.Alien
	if alienFound, found := w.aliens[alienID]; found {
		return alienFound, nil
	}

	return alien, nil
}

// AddAlien add Alien to world
func (w *WorldImpl) AddAlien(ctx context.Context, alienID int) (*types.Alien, error) {

	if _, found := w.aliens[alienID]; found {
		var alien *types.Alien
		return alien, types.ERR_DUPLICATE_ALIEN
	}

	newAlien := types.NewAlien(alienID)
	w.aliens[newAlien.AlienID] = newAlien

	return newAlien, nil
}

// MoveAlien place alien to city
func (w *WorldImpl) MoveAlien(ctx context.Context, alien *types.Alien, city *types.City) error {

	if alien == nil {
		return types.ERR_MISSING_ALIEN
	}

	if city == nil {
		return types.ERR_MISSING_CITY
	}

	alienFound, err := w.GetAlien(ctx, alien.AlienID)
	if err != nil {
		return err
	}

	if alienFound == nil {
		return types.ERR_UNKNOWN_ALIEN
	}

	cityFound, err := w.GetCity(ctx, city.Name)
	if err != nil {
		return err
	}

	if cityFound == nil {
		return types.ERR_UNKNOWN_CITY
	}

	if alien.City != nil {
		delete(w.alienInCities, alien.City)
	}

	alien.City = city
	w.alienInCities[alien.City] = alien

	return nil
}

// TrapAlien traps an alien
func (w *WorldImpl) TrapAlien(ctx context.Context, alien *types.Alien) error {

	if alien == nil {
		return types.ERR_MISSING_ALIEN
	}

	alienFound, err := w.GetAlien(ctx, alien.AlienID)
	if err != nil {
		return err
	}

	if alienFound != nil {
		delete(w.alienInCities, alienFound.City)
		w.aliens[alienFound.AlienID].IsTrapped = true
		return nil
	}

	return types.ERR_MISSING_ALIEN
}

// IsTrappedAlien check current alien is trapped
func (w *WorldImpl) IsTrappedAlien(ctx context.Context, alien *types.Alien) (bool, error) {

	if alien == nil {
		return false, types.ERR_MISSING_ALIEN
	}

	alienFound, err := w.GetAlien(ctx, alien.AlienID)
	if err != nil {
		return false, err
	}

	if alienFound != nil {
		isTrapped := w.aliens[alien.AlienID].IsTrapped
		return isTrapped, nil
	}

	return false, nil
}

// GetAlienAtCity get alien at city
func (w *WorldImpl) GetAlienAtCity(ctx context.Context, city *types.City) (*types.Alien, error) {

	var alien * types.Alien
	if city == nil {
		return alien, types.ERR_MISSING_CITY
	}

	cityFound, err := w.GetCity(ctx, city.Name)
	if err != nil {
		return alien, err
	}

	if cityFound == nil {
		return alien, types.ERR_UNKNOWN_CITY
	}

	if alienAtCity, found := w.alienInCities[city]; found {
		return alienAtCity, nil
	}

	return alien, nil
}

// GetUntrappedAliens retrieves the list of untrapped alien
func (w *WorldImpl) GetUntrappedAliens(ctx context.Context) ([]*types.Alien, error) {

	var aliens []*types.Alien
	for _, alien := range w.aliens {
		if found := w.aliens[alien.AlienID].IsTrapped; !found {
			aliens = append(aliens, alien)
		}
	}

	return aliens, nil
}