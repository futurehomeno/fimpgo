package fimptype

import (
	"encoding/json"
)

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

// EnsureInterfaces makes sure that service definition contains provided interfaces and adds them if they are missing.
func (s *Service) EnsureInterfaces(interfaces ...Interface) {
	for _, i := range interfaces {
		s.ensureInterface(i)
	}
}

// ensureInterface makes sure that service definition contains provided interface and adds it if it is missing.
func (s *Service) ensureInterface(i Interface) {
	for _, existing := range s.Interfaces {
		if existing == i {
			return
		}
	}

	s.Interfaces = append(s.Interfaces, i)
}

// PropertyStrings is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyStrings(name string) []string {
	value, ok := s.Props[name]
	if !ok {
		return nil
	}

	v, ok := value.([]string)
	if ok {
		return v
	}

	if s.cast(&v, value) {
		return v
	}

	return v
}

// PropertyString is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyString(name string) string {
	value, ok := s.Props[name]
	if !ok {
		return ""
	}

	v, ok := value.(string)
	if !ok {
		return ""
	}

	return v
}

// PropertyFloats is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyFloats(name string) []float64 {
	value, ok := s.Props[name]
	if !ok {
		return nil
	}

	v, ok := value.([]float64)
	if ok {
		return v
	}

	if s.cast(&v, value) {
		return v
	}

	return nil
}

// PropertyFloat is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyFloat(name string) (float64, bool) {
	value, ok := s.Props[name]
	if !ok {
		return 0, false
	}

	v, ok := value.(float64)
	if ok {
		return v, true
	}

	if s.cast(&v, value) {
		return v, true
	}

	return 0, false
}

// PropertyIntegers is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyIntegers(name string) []int64 {
	value, ok := s.Props[name]
	if !ok {
		return nil
	}

	v, ok := value.([]int64)
	if ok {
		return v
	}

	if s.cast(&v, value) {
		return v
	}

	return nil
}

// PropertyInteger is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyInteger(name string) (int64, bool) {
	value, ok := s.Props[name]
	if !ok {
		return 0, false
	}

	v, ok := value.(int64)
	if ok {
		return v, true
	}

	if s.cast(&v, value) {
		return v, true
	}

	return 0, false
}

// PropertyBool is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyBool(name string) bool {
	value, ok := s.Props[name]
	if !ok {
		return false
	}

	v, ok := value.(bool)
	if ok {
		return v
	}

	if s.cast(&v, value) {
		return v
	}

	return false
}

// PropertyObject is a helper that extracts property settings out of the service specification.
func (s *Service) PropertyObject(name string, object interface{}) bool {
	value, ok := s.Props[name]
	if !ok {
		return false
	}

	return s.cast(object, value)
}

// cast is a helper allowing simple casting of interfaces to destination type using marshalling and unmarshalling in the process.
func (s *Service) cast(dst, src interface{}) bool {
	b, err := json.Marshal(src)
	if err != nil {
		return false
	}

	err = json.Unmarshal(b, dst)

	return err == nil
}

// Constants defining type of interface.
const (
	TypeIn  = "in"
	TypeOut = "out"
)

// Interface represents a supported communication interface with the service.
type Interface struct {
	Type      string `json:"intf_t"`
	MsgType   string `json:"msg_t"`
	ValueType string `json:"val_t"`
	Version   string `json:"ver"`
}
