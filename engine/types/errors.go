package types

import (
	"fmt"
)

var (
	ERR_PARSE_CITY_DEFINITION error = fmt.Errorf("error parsing the city definition")

	ERR_EMPTY_CITY_NAME error = fmt.Errorf("city name is Empty")

	ERR_DUPLICATE_CITY error = fmt.Errorf("duplicate city name exists")

	ERR_MISSING_CITY error = fmt.Errorf("city is missing")

	ERR_UNKNOWN_CITY error = fmt.Errorf("city is unknown")

	ERR_LINK_SAME_CITY error = fmt.Errorf("no possible link between same city")

	ERR_DUPLICATE_ALIEN error = fmt.Errorf("duplicate alien not allowed")

	ERR_MISSING_ALIEN error = fmt.Errorf("alien is missing")

	ERR_UNKNOWN_ALIEN error = fmt.Errorf("alien is unkown")

	ERR_UNKNOWN_DIRECTION error = fmt.Errorf("unknown direction provided")
	
	ERR_ALREADY_EXISTS_LINK error = fmt.Errorf("a link already exists between the two cities")

	ERR_RANDOM_OUT_OF_BOUNDS  error = fmt.Errorf("random input out of bounds")

	ERR_CONTEXT_CANCELLED  error = fmt.Errorf("the context was cancelled")

)