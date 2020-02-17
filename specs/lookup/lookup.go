package lookup

import (
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

// Reference represents a property reference
type Reference interface {
	GetPath() string
	GetDefault() interface{}
	GetType() types.Type
	GetObject() specs.Object
}

// ReferenceMap holds the resource references and their representing parameter map
type ReferenceMap map[string]PathLookup

// PathLookup represents a lookup method that returns the property available on the given path
type PathLookup func(path string) Reference

// GetFlow attempts to find the given flow inside the given manifest
func GetFlow(manifest specs.Manifest, name string) *specs.Flow {
	for _, flow := range manifest.Flows {
		if flow.Name == name {
			return flow
		}
	}

	return nil
}

// GetDefaultProp returns the default resource for the given resource
func GetDefaultProp(resource string) string {
	if resource == specs.InputResource {
		return specs.ResourceRequest
	}

	return specs.ResourceResponse
}

// GetAvailableResources fetches the available resources able to be referenced
// untill the given breakpoint (call.Name) has been reached.
func GetAvailableResources(flow specs.FlowManager, breakpoint string) map[string]ReferenceMap {
	references := make(map[string]ReferenceMap, len(flow.GetCalls())+1)

	if flow.GetInput() != nil {
		references[specs.InputResource] = ReferenceMap{
			specs.ResourceRequest:       ParameterMapLookup(flow.GetInput()),
			specs.ResourceRequestHeader: HeaderLookup(flow.GetInput().Header),
		}
	}

	for _, call := range flow.GetCalls() {
		if call.Name == breakpoint {
			break
		}

		resources := ReferenceMap{}

		if call.Request != nil {
			resources[specs.ResourceRequest] = ParameterMapLookup(call.Request)
			resources[specs.ResourceRequestHeader] = HeaderLookup(call.Request.Header)
		}

		if call.Response != nil {
			resources[specs.ResourceResponse] = ParameterMapLookup(call.Response)
		}

		references[call.Name] = resources
	}

	return references
}

// GetResourceReference attempts to return the resource reference property
func GetResourceReference(reference *specs.PropertyReference, references map[string]ReferenceMap) Reference {
	resources := specs.SplitPath(reference.Resource)

	target := resources[0]
	prop := GetDefaultProp(target)

	if len(resources) > 1 {
		prop = specs.JoinPath(resources[1:]...)
	}

	for resource, refs := range references {
		if resource != target {
			continue
		}

		return GetReference(reference.Path, prop, refs)
	}

	return nil
}

// GetReference attempts to lookup and return the available property on the given path
func GetReference(path string, prop string, references ReferenceMap) Reference {
	lookup, has := references[prop]
	if !has {
		return nil
	}

	return lookup(path)
}

// HeaderLookup attempts to lookup the given path inside the header
func HeaderLookup(header specs.Header) PathLookup {
	return func(path string) Reference {
		for key, header := range header {
			if key == path {
				return header
			}
		}

		return nil
	}
}

// ParameterMapLookup attempts to lookup the given path inside the params collection
func ParameterMapLookup(params specs.Object) PathLookup {
	return func(path string) Reference {
		for _, param := range params.GetProperties() {
			if param.GetPath() == path {
				return param
			}
		}

		for _, nested := range params.GetNestedProperties() {
			if nested.GetPath() == path {
				return nested
			}

			prop := ParameterMapLookup(nested)(path)
			if prop != nil {
				return prop
			}
		}

		for _, repeated := range params.GetRepeatedProperties() {
			if repeated.GetPath() == path {
				return repeated
			}

			prop := ParameterMapLookup(repeated)(path)
			if prop != nil {
				return prop
			}
		}

		return nil
	}
}
