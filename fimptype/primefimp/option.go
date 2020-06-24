package primefimp

type apiClientConfig struct {
	cloudService string
}

type Option interface {
	apply(*apiClientConfig)
}

type cloudServiceOption string

func (cso cloudServiceOption) apply(config *apiClientConfig) {
	config.cloudService = string(cso)
}

func WithCloudService(service string) Option {
	return cloudServiceOption(service)
}
