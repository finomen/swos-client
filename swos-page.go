package swos_client

type swOsPage[I any] = interface {
	url() string
	load(I) error
}

type writebleSwOsPage[O any, I any] = interface {
	store() O
	swOsPage[I]
}
