package commercetools

type Tax struct {
	Country         string
	Amount          float64
	Name            string
	IncludedInPrice bool `yaml:"included_in_price" default:"true"`
}
