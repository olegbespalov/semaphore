package lookup

import (
	"testing"

	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/types"
)

func TestGetFlow(t *testing.T) {
	manifest := specs.FlowsManifest{
		Flows: []*specs.Flow{
			{
				Name: "first",
			},
			{
				Name: "second",
			},
		},
	}

	tests := map[string]*specs.Flow{
		"second":  manifest.Flows[1],
		"first":   manifest.Flows[0],
		"unknown": nil,
	}

	for input, expected := range tests {
		result := GetFlow(manifest, input)
		if result != expected {
			t.Errorf("unexpected result %+v, expected %+v", result, expected)
		}
	}
}

func TestGetDefaultProp(t *testing.T) {
	tests := map[string]string{
		template.InputResource: template.ResourceRequest,
		"unknown":              template.ResourceResponse,
	}

	for input, expected := range tests {
		result := GetDefaultProp(input)
		if result != expected {
			t.Errorf("unexpected result %s, expected %s", result, expected)
		}
	}
}

func TestGetNextResource(t *testing.T) {
	type test struct {
		breakpoint string
		expected   string
		manager    *specs.Flow
	}

	tests := map[string]*test{
		"first": {
			breakpoint: "first",
			expected:   "second",
			manager: &specs.Flow{
				Nodes: []*specs.Node{
					{
						Name: "first",
					},
					{
						Name: "second",
					},
				},
			},
		},
		"second": {
			breakpoint: "second",
			expected:   "third",
			manager: &specs.Flow{
				Nodes: []*specs.Node{
					{
						Name: "first",
					},
					{
						Name: "second",
					},
					{
						Name: "third",
					},
				},
			},
		},
		"output": {
			breakpoint: "last",
			expected:   template.OutputResource,
			manager: &specs.Flow{
				Nodes: []*specs.Node{
					{
						Name: "first",
					},
					{
						Name: "last",
					},
				},
			},
		},
		"unkown": {
			breakpoint: "unkown",
			expected:   "unkown",
			manager: &specs.Flow{
				Nodes: []*specs.Node{
					{
						Name: "first",
					},
					{
						Name: "second",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := GetNextResource(test.manager, test.breakpoint)
			if result != test.expected {
				t.Fatalf("unexpected result '%s', expected '%s'", result, test.expected)
			}
		})
	}
}

func NewInputMockProperty() *specs.Property {
	return &specs.Property{
		Path:  "",
		Type:  types.Message,
		Label: labels.Optional,
		Nested: map[string]*specs.Property{
			"message": {
				Path:    "message",
				Default: "hello world",
				Type:    types.String,
				Label:   labels.Optional,
			},
			"active": {
				Path:    "active",
				Default: false,
				Type:    types.Bool,
				Label:   labels.Optional,
			},
			"nested": {
				Path:  "nested",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*specs.Property{
					"message": {
						Path:    "nested.message",
						Default: "hello world",
						Type:    types.String,
						Label:   labels.Optional,
					},
					"active": {
						Path:    "nested.active",
						Default: false,
						Type:    types.Bool,
						Label:   labels.Optional,
					},
					"nested": {
						Path:  "nested.nested",
						Type:  types.Message,
						Label: labels.Optional,
						Nested: map[string]*specs.Property{
							"message": {
								Path:    "nested.nested.message",
								Default: "hello world",
								Type:    types.String,
								Label:   labels.Optional,
							},
							"active": {
								Path:    "nested.nested.active",
								Default: false,
								Type:    types.Bool,
								Label:   labels.Optional,
							},
						},
					},
					"repeated": {
						Path:  "nested.repeated",
						Type:  types.Message,
						Label: labels.Repeated,
						Nested: map[string]*specs.Property{
							"message": {
								Path:    "nested.repeated.message",
								Default: "hello world",
								Type:    types.String,
								Label:   labels.Optional,
							},
							"active": {
								Path:    "nested.repeated.active",
								Default: false,
								Type:    types.Bool,
								Label:   labels.Optional,
							},
						},
					},
				},
			},
			"repeated": {
				Path:  "repeated",
				Type:  types.Message,
				Label: labels.Repeated,
				Nested: map[string]*specs.Property{
					"message": {
						Path:    "message",
						Default: "hello world",
						Type:    types.String,
						Label:   labels.Optional,
					},
					"active": {
						Path:    "active",
						Default: false,
						Type:    types.Bool,
						Label:   labels.Optional,
					},
					"repeated": {
						Path:  "repeated.repeated",
						Type:  types.Message,
						Label: labels.Repeated,
						Nested: map[string]*specs.Property{
							"repeated.message": {
								Path:    "message",
								Default: "hello world",
								Type:    types.String,
								Label:   labels.Optional,
							},
							"repeated.active": {
								Path:    "active",
								Default: false,
								Type:    types.Bool,
								Label:   labels.Optional,
							},
						},
					},
				},
			},
		},
	}
}

func NewResultMockProperty() *specs.Property {
	return &specs.Property{
		Path:  "",
		Type:  types.Message,
		Label: labels.Optional,
		Nested: map[string]*specs.Property{
			"result": {
				Path:    "result",
				Default: "hello world",
				Type:    types.String,
				Label:   labels.Optional,
			},
			"active": {
				Path:    "active",
				Default: false,
				Type:    types.Bool,
				Label:   labels.Optional,
			},
			"nested": {
				Path:  "nested",
				Type:  types.Message,
				Label: labels.Optional,
				Nested: map[string]*specs.Property{
					"result": {
						Path:    "nested.result",
						Default: "hello world",
						Type:    types.String,
						Label:   labels.Optional,
					},
					"active": {
						Path:    "nested.active",
						Default: false,
						Type:    types.Bool,
						Label:   labels.Optional,
					},
					"nested": {
						Path:  "nested.nested",
						Type:  types.Message,
						Label: labels.Optional,
						Nested: map[string]*specs.Property{
							"result": {
								Path:    "nested.nested.result",
								Default: "hello world",
								Type:    types.String,
								Label:   labels.Optional,
							},
							"active": {
								Path:    "nested.nested.active",
								Default: false,
								Type:    types.Bool,
								Label:   labels.Optional,
							},
						},
					},
					"repeated": {
						Path:  "nested.repeated",
						Type:  types.Message,
						Label: labels.Repeated,
						Nested: map[string]*specs.Property{
							"result": {
								Path:    "nested.repeated.result",
								Default: "hello world",
								Type:    types.String,
								Label:   labels.Optional,
							},
							"active": {
								Path:    "nested.repeated.active",
								Default: false,
								Type:    types.Bool,
								Label:   labels.Optional,
							},
						},
					},
				},
			},
			"repeated": {
				Path:  "repeated",
				Type:  types.Message,
				Label: labels.Repeated,
				Nested: map[string]*specs.Property{
					"result": {
						Path:    "message",
						Default: "hello world",
						Type:    types.String,
						Label:   labels.Optional,
					},
					"active": {
						Path:    "active",
						Default: false,
						Type:    types.Bool,
						Label:   labels.Optional,
					},
					"repeated": {
						Path:  "repeated.repeated",
						Type:  types.Message,
						Label: labels.Repeated,
						Nested: map[string]*specs.Property{
							"result": {
								Path:    "repeated.repeated.result",
								Default: "hello world",
								Type:    types.String,
								Label:   labels.Optional,
							},
							"active": {
								Path:    "repeated.repeated.active",
								Default: false,
								Type:    types.Bool,
								Label:   labels.Optional,
							},
						},
					},
				},
			},
		},
	}
}

func NewMockCall(name string) *specs.Node {
	return &specs.Node{
		Name: name,
		Call: &specs.Call{
			Request: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:    "cookie",
						Default: "mnomnom",
						Type:    types.String,
						Label:   labels.Optional,
					},
				},
				Property: NewInputMockProperty(),
			},
			Response: &specs.ParameterMap{
				Header: specs.Header{
					"cookie": &specs.Property{
						Path:    "cookie",
						Default: "mnomnom",
						Type:    types.String,
						Label:   labels.Optional,
					},
				},
				Property: NewResultMockProperty(),
			},
		},
	}
}

func NewMockFlow(name string) *specs.Flow {
	return &specs.Flow{
		Name: name,
		Input: &specs.ParameterMap{
			Property: NewInputMockProperty(),
		},
		Nodes: []*specs.Node{
			NewMockCall("first"),
			NewMockCall("second"),
			NewMockCall("third"),
		},
		Output: &specs.ParameterMap{
			Property: NewResultMockProperty(),
		},
	}
}

func TestGetAvailableResources(t *testing.T) {
	tests := map[string]func() ([]string, map[string]ReferenceMap){
		"input and first": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, "input", "first", "second"}

			result := GetAvailableResources(flow, "second")
			return expected, result
		},
		"input": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, "input", "first"}

			result := GetAvailableResources(flow, "first")
			return expected, result
		},
		"output": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, "input", "first", "second", "third"}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
		"stack lookup": func() ([]string, map[string]ReferenceMap) {
			flow := NewMockFlow("first")
			expected := []string{template.StackResource, "input", "first", "second", "third"}

			flow.Nodes[0].Call.Request.Stack = map[string]*specs.Property{
				"ref": {
					Path: "ref",
				},
			}

			flow.Nodes[0].Call.Response.Stack = map[string]*specs.Property{
				"ref": {
					Path: "ref",
				},
			}

			result := GetAvailableResources(flow, "output")
			return expected, result
		},
	}

	for key, test := range tests {
		t.Run(key, func(t *testing.T) {
			expected, result := test()
			if len(expected) != len(result) {
				t.Errorf("unexpected result %+v, expected %+v", result, expected)
			}

			for _, resource := range expected {
				_, has := result[resource]
				if !has {
					t.Errorf("expected resource not found %s, %+v", resource, result)
				}
			}
		})
	}
}

func TestSkipMissingParameters(t *testing.T) {
	flow := NewMockFlow("first")

	flow.Input = nil
	flow.Nodes[0].Call.Request = nil
	flow.Nodes[1].Call.Response = nil

	type test struct {
		name  string
		empty []string
	}

	tests := []test{
		{
			name: template.StackResource,
		},
		{
			name:  "first",
			empty: []string{template.ResourceRequest, template.ResourceHeader},
		},
		{
			name:  "second",
			empty: []string{template.ResourceRequest, template.ResourceHeader},
		},
		{
			name: "third",
		},
	}

	result := GetAvailableResources(flow, "output")

	if len(result) != len(tests) {
		t.Fatalf("unexpected result (%d, %d) %+v, expected %+v", len(result), len(tests), result, tests)
	}

	for _, test := range tests {
		resource, has := result[test.name]
		if !has {
			t.Fatalf("expected resource not found %s, %+v", test.name, result)
		}

	check:
		for key, val := range resource {
			if val == nil {
				for _, known := range test.empty {
					if key == known {
						continue check
					}
				}

				t.Fatalf("unknown empty resource %s", key)
			}
		}
	}
}

func NewPropertyReference(resource string, path string) *specs.PropertyReference {
	return &specs.PropertyReference{
		Resource: resource,
		Path:     path,
	}
}

func TestGetResourceReference(t *testing.T) {
	flow := NewMockFlow("first")
	references := GetAvailableResources(flow, "output")
	breakpoint := "first"

	tests := map[*specs.PropertyReference]*specs.Property{
		NewPropertyReference("input", "message"):                           flow.Input.Property.Nested["message"],
		NewPropertyReference("first", "result"):                            flow.Nodes[0].Call.Response.Property.Nested["result"],
		NewPropertyReference("", "result"):                                 flow.Nodes[0].Call.Response.Property.Nested["result"],
		NewPropertyReference("first.response", "result"):                   flow.Nodes[0].Call.Response.Property.Nested["result"],
		NewPropertyReference("first.request", "message"):                   flow.Nodes[0].Call.Request.Property.Nested["message"],
		NewPropertyReference("first.request", "message"):                   flow.Nodes[0].Call.Request.Property.Nested["message"],
		NewPropertyReference("first.header", "cookie"):                     flow.Nodes[0].Call.Response.Header["cookie"],
		NewPropertyReference("first.request", "nested"):                    flow.Nodes[0].Call.Request.Property.Nested["nested"],
		NewPropertyReference("first.request", "nested.message"):            flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["message"],
		NewPropertyReference("first", "nested.result"):                     flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["result"],
		NewPropertyReference("first", "nested.nested.result"):              flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["nested"].Nested["result"],
		NewPropertyReference("first.response", "nested.nested.result"):     flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["nested"].Nested["result"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["repeated"].Nested["message"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.response", "repeated.repeated.result"): flow.Nodes[0].Call.Response.Property.Nested["repeated"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.response", "nested.nested.result"):     flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["nested"].Nested["result"],
		NewPropertyReference("first.response", "nested.repeated.result"):   flow.Nodes[0].Call.Response.Property.Nested["nested"].Nested["repeated"].Nested["result"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["repeated"].Nested["message"],
		NewPropertyReference("first.request", "nested.nested.message"):     flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["nested"].Nested["message"],
		NewPropertyReference("first.request", "nested.repeated.message"):   flow.Nodes[0].Call.Request.Property.Nested["nested"].Nested["repeated"].Nested["message"],
	}

	for input, expected := range tests {
		t.Run(input.String(), func(t *testing.T) {
			result := GetResourceReference(input, references, breakpoint)
			if result == nil {
				t.Fatalf("unexpected empty result on lookup '%s', expected '%+v'", input, expected)
			}

			if result.Path != expected.Path {
				t.Fatalf("unexpected result '%+v', expected '%+v'", result, expected)
			}
		})
	}
}

func TestGetUnkownResourceReference(t *testing.T) {
	flow := NewMockFlow("first")
	references := GetAvailableResources(flow, "output")
	breakpoint := "first"

	tests := map[string]*specs.PropertyReference{
		"unkown": NewPropertyReference("unkown", "unkown"),
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := GetResourceReference(test, references, breakpoint)
			if result != nil {
				t.Fatalf("unexpected result")
			}
		})
	}
}

func TestHeaderLookup(t *testing.T) {
	type test struct {
		path   string
		header specs.Header
	}

	tests := map[string]*test{
		"simple": {
			path: "key",
			header: specs.Header{
				"key": &specs.Property{
					Path: "key",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := HeaderLookup(test.header)
			prop := resolver(test.path)
			if prop == nil {
				t.Fatalf("unexpected result expected a prop to be returned for '%s'", test.path)
			}
		})
	}
}

func TestUnkownHeaderLookup(t *testing.T) {
	type test struct {
		path   string
		header specs.Header
	}

	tests := map[string]*test{
		"simple": {
			path:   "key",
			header: specs.Header{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := HeaderLookup(test.header)
			prop := resolver(test.path)
			if prop != nil {
				t.Fatalf("unexpected result expected nil to be returned for '%s'", test.path)
			}
		})
	}
}

func TestPropertyLookup(t *testing.T) {
	type test struct {
		path  string
		param *specs.Property
	}

	tests := map[string]*test{
		"self reference": {
			path: ".",
			param: &specs.Property{
				Path: "key",
			},
		},
		"simple": {
			path: "key",
			param: &specs.Property{
				Path: "key",
			},
		},
		"nested": {
			path: "key.nested",
			param: &specs.Property{
				Path: "key",
				Nested: map[string]*specs.Property{
					"nested": {
						Path: "key.nested",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lookup := PropertyLookup(test.param)
			result := lookup(test.path)
			if result == nil {
				t.Fatal("unexpected empty result")
			}
		})
	}
}

func TestResolveSelfReference(t *testing.T) {
	type test struct {
		path     string
		expected string
		resource string
	}

	tests := map[string]*test{
		"self reference": {
			path:     ".request",
			resource: "input",
			expected: "input.request",
		},
		"resource reference": {
			path:     "input.request",
			resource: "first",
			expected: "input.request",
		},
		"broken path": {
			path:     "input.",
			resource: "first",
			expected: "input.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ResolveSelfReference(test.path, test.resource)
			if result != test.expected {
				t.Fatalf("unexpected result '%s', expected '%s'", result, test.expected)
			}
		})
	}
}

func TestGetReference(t *testing.T) {
	path := "key"
	prop := "input.request"

	references := ReferenceMap{
		prop: PropertyLookup(&specs.Property{Path: path}),
	}

	result := GetReference(path, prop, references)
	if result == nil {
		t.Fatal("unexpected empty result")
	}
}

func TestUnkownReference(t *testing.T) {
	path := "key"
	prop := "input.request"

	references := ReferenceMap{
		prop: PropertyLookup(&specs.Property{Path: path}),
	}

	result := GetReference(path, "unkown", references)
	if result != nil {
		t.Fatal("unexpected result")
	}
}
