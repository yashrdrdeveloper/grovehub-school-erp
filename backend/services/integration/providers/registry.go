package providers

import (
 "errors"
 "sort"
 "strings"
)

var ErrUnknownProvider = errors.New("unknown integration provider")
var ErrDuplicateProvider = errors.New("duplicate integration provider")

type Registry struct { providers map[string]Provider }

func NewRegistry(items ...Provider) (*Registry, error) {
 registry := &Registry{providers: make(map[string]Provider, len(items))}
 for _, item := range items {
  descriptor := item.Descriptor()
  id := strings.ToUpper(strings.TrimSpace(descriptor.ID))
  if id == "" { return nil, ErrUnknownProvider }
  if _, exists := registry.providers[id]; exists { return nil, ErrDuplicateProvider }
  registry.providers[id] = item
 }
 return registry, nil
}
func (r *Registry) Find(id string) (Provider, error) {
 provider, ok := r.providers[strings.ToUpper(strings.TrimSpace(id))]
 if !ok { return nil, ErrUnknownProvider }
 return provider, nil
}
func (r *Registry) List() []Descriptor {
 result := make([]Descriptor, 0, len(r.providers))
 for _, item := range r.providers {
  d := item.Descriptor()
  d.Capabilities = append([]Capability(nil), d.Capabilities...)
  result = append(result, d)
 }
 sort.Slice(result, func(i,j int) bool { return result[i].ID < result[j].ID })
 return result
}
