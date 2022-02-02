package fimptype

// Service represents a specification of the service supported by the thing.
type Service struct {
	Name             string                 `json:"name" storm:"index"`
	Alias            string                 `json:"alias"`
	Address          string                 `json:"address"`
	Enabled          bool                   `json:"enabled"`
	Groups           []string               `json:"groups"`
	Props            map[string]interface{} `json:"props"`
	Tags             []string               `json:"tags"`
	PropSetReference string                 `json:"prop_set_ref"`
	Interfaces       []Interface            `json:"interfaces"`
}

// PropertyStrings is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyStrings(name string) []string {
	value, ok := s.Props[name]
	if !ok {
		return nil
	}

	values, ok := value.([]string)
	if !ok {
		return nil
	}

	return values
}

// Interface represents a supported communication interface with the service.
type Interface struct {
	Type      string `json:"intf_t"`
	MsgType   string `json:"msg_t"`
	ValueType string `json:"val_t"`
	Version   string `json:"ver"`
}
