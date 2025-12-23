package primefimp

type (
	apiClientConfig struct {
		cloudService string
		globalPrefix string
	}
)

type Option interface {
	apply(*apiClientConfig)
}

type cloudServiceOption string

func (cso cloudServiceOption) apply(config *apiClientConfig) {
	config.cloudService = string(cso)
}

type globalPrefixOption string

func (cso globalPrefixOption) apply(config *apiClientConfig) {
	config.globalPrefix = string(cso)
}

func WithCloudService(service string) Option {
	return cloudServiceOption(service)
}

func WithGlobalPrefix(prefix string) Option {
	return globalPrefixOption(prefix)
}
