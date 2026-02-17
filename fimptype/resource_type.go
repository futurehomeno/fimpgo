package fimptype

type ResourceTypeT string

const (
	ResourceTypeDevice    ResourceTypeT = "dev"
	ResourceTypeApp       ResourceTypeT = "app"
	ResourceTypeAdapter   ResourceTypeT = "ad"
	ResourceTypeCloud     ResourceTypeT = "cloud"
	ResourceTypeDiscovery ResourceTypeT = "discovery"
	ResourceTypeLocation  ResourceTypeT = "loc"
	ResourceTypeUnknown   ResourceTypeT = ""
)

func (rn ResourceTypeT) Str() string {
	return string(rn)
}
