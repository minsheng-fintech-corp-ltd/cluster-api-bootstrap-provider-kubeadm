package ignition

type ServiceUnit struct {
	Content string
	Dropins []Dropin
	Enabled bool
	Name    string
}

type Dropin struct {
	Name    string
	Content string
}
