package references

import (
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/lookup"
	"github.com/jexia/maestro/pkg/specs/template"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/sirupsen/logrus"
)

// DefineManifest checks and defines the types for the given manifest
func DefineManifest(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest) (err error) {
	ctx.Logger(logger.Core).Info("Defining manifest types")

	for _, flow := range flows.Flows {
		err := DefineFlow(ctx, services, schema, flows, flow)
		if err != nil {
			return err
		}
	}

	for _, proxy := range flows.Proxy {
		err := DefineProxy(ctx, services, schema, flows, proxy)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineProxy checks and defines the types for the given proxy
func DefineProxy(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, proxy *specs.Proxy) (err error) {
	ctx.Logger(logger.Core).WithField("proxy", proxy.GetName()).Info("Defining proxy flow types")

	if proxy.Input != nil && proxy.Input.Schema != "" {
		input := schema.GetProperty(proxy.Input.Schema)
		if input == nil {
			return trace.New(trace.WithMessage("undefined object '%s' in schema collection", proxy.Input.Schema))
		}

		proxy.Input = ToParameterMap(proxy.Input, "", input)
	}

	if proxy.Error != nil {
		prop := schema.GetProperty(proxy.Error.Schema)
		if prop == nil {
			return trace.New(trace.WithMessage("undefined error object '%s' in schema collection", proxy.Error.Schema))
		}

		proxy.Error = ToError(prop, proxy.Error)
	}

	for _, node := range proxy.Nodes {
		if node.Call != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Call, proxy)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Rollback, proxy)
			if err != nil {
				return err
			}
		}

		if node.OnError != nil {
			err = DefineError(ctx, services, proxy, node, node.OnError)
			if err != nil {
				return err
			}
		}
	}

	if proxy.Forward != nil {
		for _, header := range proxy.Forward.Request.Header {
			err = DefineProperty(ctx, nil, header, proxy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DefineFlow defines the types for the given flow
func DefineFlow(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, flow *specs.Flow) (err error) {
	ctx.Logger(logger.Core).WithField("flow", flow.GetName()).Info("Defining flow types")

	if flow.Input != nil {
		input := schema.GetProperty(flow.Input.Schema)
		if input == nil {
			return trace.New(trace.WithMessage("undefined object '%s' in schema collection", flow.Input.Schema))
		}

		flow.Input = ToParameterMap(flow.Input, "", input)
	}

	if flow.Error != nil {
		prop := schema.GetProperty(flow.Error.Schema)
		if prop == nil {
			return trace.New(trace.WithMessage("undefined error object '%s' in schema collection", flow.Error.Schema))
		}

		flow.Error = ToError(prop, flow.Error)
	}

	for _, node := range flow.Nodes {
		if node.Call != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Call, flow)
			if err != nil {
				return err
			}
		}

		if node.Rollback != nil {
			err = DefineCall(ctx, services, schema, flows, node, node.Rollback, flow)
			if err != nil {
				return err
			}
		}

		if node.OnError != nil {
			err = DefineError(ctx, services, flow, node, node.OnError)
			if err != nil {
				return err
			}
		}
	}

	if flow.Output != nil {
		err = DefineParameterMap(ctx, nil, flow.Output, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineError defined the types for the given error
func DefineError(ctx instance.Context, services *specs.ServicesManifest, flow specs.FlowResourceManager, node *specs.Node, onError *specs.OnError) error {
	for _, param := range onError.Params {
		reference, err := LookupReference(ctx, node, node.Name, param, flow)
		if err != nil {
			return err
		}

		param.Property = reference
	}

	return nil
}

// DefineCall defineds the types for the specs call
func DefineCall(ctx instance.Context, services *specs.ServicesManifest, schema *specs.SchemaManifest, flows *specs.FlowsManifest, node *specs.Node, call *specs.Call, flow specs.FlowResourceManager) (err error) {
	if call.Request != nil {
		err = DefineParameterMap(ctx, node, call.Request, flow)
		if err != nil {
			return err
		}
	}

	if call.Method != "" {
		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"call":    node.Name,
			"method":  call.Method,
			"service": call.Service,
		}).Info("Defining call types")

		service := services.GetService(call.Service)
		if service == nil {
			return trace.New(trace.WithMessage("undefined service '%s' in flow '%s'", call.Service, flow.GetName()))
		}

		method := service.GetMethod(call.Method)
		if method == nil {
			return trace.New(trace.WithMessage("undefined method '%s' in flow '%s'", call.Method, flow.GetName()))
		}

		output := schema.GetProperty(method.Output)
		if output == nil {
			return trace.New(trace.WithMessage("undefined method output property '%s' in flow '%s'", method.Output, flow.GetName()))
		}

		call.Descriptor = method
		call.Response = ToParameterMap(nil, "", output)
	}

	if call.Response != nil {
		err = DefineParameterMap(ctx, node, call.Response, flow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefineParameterMap defines the types for the given parameter map
func DefineParameterMap(ctx instance.Context, node *specs.Node, params *specs.ParameterMap, flow specs.FlowResourceManager) (err error) {
	if params.Property == nil {
		return nil
	}

	for _, header := range params.Header {
		err = DefineProperty(ctx, node, header, flow)
		if err != nil {
			return err
		}
	}

	err = DefineParams(ctx, node, params.Params, flow)
	if err != nil {
		return err
	}

	err = DefineProperty(ctx, node, params.Property, flow)
	if err != nil {
		return err
	}

	return nil
}

// DefineParams defines all types inside the given params
func DefineParams(ctx instance.Context, node *specs.Node, params map[string]*specs.PropertyReference, flow specs.FlowResourceManager) error {
	if params == nil {
		return nil
	}

	for _, param := range params {
		reference, err := LookupReference(ctx, node, node.Name, param, flow)
		if err != nil {
			return err
		}

		param.Property = reference
	}

	return nil
}

// DefineProperty defines the given property type.
// If any object is references it has to be fixed afterwards and moved into the correct dataset
func DefineProperty(ctx instance.Context, node *specs.Node, property *specs.Property, flow specs.FlowResourceManager) error {
	if len(property.Nested) > 0 {
		for _, nested := range property.Nested {
			err := DefineProperty(ctx, node, nested, flow)
			if err != nil {
				return err
			}
		}
	}

	if property.Reference == nil {
		return nil
	}

	breakpoint := template.OutputResource
	if node != nil {
		breakpoint = node.Name

		if node.Rollback != nil && property != nil {
			rollback := node.Rollback.Request.Property
			if InsideProperty(rollback, property) {
				breakpoint = lookup.GetNextResource(flow, breakpoint)
			}
		}
	}

	reference, err := LookupReference(ctx, node, breakpoint, property.Reference, flow)
	if err != nil {
		return trace.New(trace.WithExpression(property.Expr), trace.WithMessage("undefined resource '%s' in '%s.%s.%s'", property.Reference, flow.GetName(), breakpoint, property.Path))
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"reference": property.Reference,
		"name":      reference.Name,
		"path":      reference.Path,
	}).Debug("References lookup result")

	property.Type = reference.Type
	property.Label = reference.Label
	property.Default = reference.Default
	property.Reference.Property = reference

	if reference.Enum != nil {
		property.Enum = reference.Enum
	}

	return nil
}

// LookupReference looks up the given reference
func LookupReference(ctx instance.Context, node *specs.Node, breakpoint string, reference *specs.PropertyReference, flow specs.FlowResourceManager) (*specs.Property, error) {
	reference.Resource = lookup.ResolveSelfReference(reference.Resource, breakpoint)

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"breakpoint": breakpoint,
		"reference":  reference,
	}).Debug("Lookup references until breakpoint")

	references := lookup.GetAvailableResources(flow, breakpoint)
	result := lookup.GetResourceReference(reference, references, breakpoint)
	if result == nil {
		return nil, trace.New(trace.WithMessage("undefined resource '%s' in '%s'.'%s'", reference, flow.GetName(), breakpoint))
	}

	ctx.Logger(logger.Core).WithFields(logrus.Fields{
		"breakpoint": breakpoint,
		"reference":  result,
	}).Debug("Lookup references result")

	return result, nil
}

// InsideProperty checks whether the given property is insde the source property
func InsideProperty(source *specs.Property, target *specs.Property) bool {
	if source == target {
		return true
	}

	if len(source.Nested) > 0 {
		for _, nested := range source.Nested {
			is := InsideProperty(nested, target)
			if is {
				return is
			}
		}
	}

	return false
}
