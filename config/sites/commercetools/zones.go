package commercetools

import "github.com/labd/mach-composer/internal/config/sites/commercetools/zones"

type Zone struct {
	Name        string
	Description string
	Locations   []zones.Location
}
