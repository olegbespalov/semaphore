package xml

import (
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

var (
	propString = &specs.Property{
		Name:  "string",
		Path:  "string",
		Label: labels.Required,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.String,
			},
		},
	}

	propInteger = &specs.Property{
		Name:  "integer",
		Path:  "integer",
		Label: labels.Required,
		Template: specs.Template{
			Scalar: &specs.Scalar{
				Type: types.Int32,
			},
		},
	}

	propEnum = &specs.Property{
		Name:  "status",
		Path:  "status",
		Label: labels.Required,
		Template: specs.Template{
			Enum: enum,
		},
	}

	SchemaEnum = &specs.ParameterMap{
		Property: propEnum,
	}

	SchemaScalar = &specs.ParameterMap{
		Property: propInteger,
	}

	SchemaArray = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "array",
			Path:  "array",
			Label: labels.Optional,
			Template: specs.Template{
				Repeated: specs.Repeated{
					propString.Template,
				},
			},
		},
	}

	// TODO: array in message

	SchemaObject = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "root",
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"status": func() *specs.Property {
						var clone = propEnum.Clone()
						clone.Position = 1
						clone.Path = "root." + clone.Path

						return clone
					}(),
					"integer": func() *specs.Property {
						var clone = propInteger.Clone()
						clone.Position = 2
						clone.Path = "root." + clone.Path

						return clone
					}(),
				},
			},
		},
	}

	SchemaObjectNested = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "root",
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"nested": {
						Position: 1,
						Name:     "nested",
						Path:     "root.nested",
						Label:    labels.Optional,
						Template: specs.Template{
							Message: specs.Message{
								"status": func() *specs.Property {
									var clone = propEnum.Clone()
									clone.Position = 1
									clone.Path = "root.nested." + clone.Path

									return clone
								}(),
								"integer": func() *specs.Property {
									var clone = propInteger.Clone()
									clone.Position = 2
									clone.Path = "root.nested." + clone.Path

									return clone
								}(),
							},
						},
					},
					"string": func() *specs.Property {
						var clone = propString.Clone()
						clone.Position = 2
						clone.Path = "root." + clone.Path

						return clone
					}(),
				},
			},
		},
	}

	SchemaComplexObject = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "root",
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"bad_label": {
						Position: 1,
						Name:     "bad_label",
						Path:     "bad_label",
						Label:    "unknown",
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"numeric": {
						Position: 3,
						Name:     "numeric",
						Path:     "numeric",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "numeric",
							},
							Scalar: &specs.Scalar{
								Type: types.Int32,
							},
						},
					},
					"message": {
						Position: 4,
						Name:     "message",
						Path:     "message",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "message",
							},
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"another_message": {
						Position: 5,
						Name:     "another_message",
						Path:     "another_message",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "another_message",
							},
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"status": {
						Position: 6,
						Name:     "status",
						Path:     "status",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "status",
							},
							Enum: enum,
						},
					},
					"another_status": {
						Position: 7,
						Name:     "another_status",
						Path:     "another_status",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "another_status",
							},
							Enum: enum,
						},
					},
					"nested": {
						Position: 8,
						Name:     "nested",
						Path:     "nested",
						Label:    labels.Optional,
						Template: specs.Template{
							Message: specs.Message{
								"first": {
									Position: 1,
									Name:     "first",
									Path:     "nested.first",
									Label:    labels.Optional,
									Template: specs.Template{
										Reference: &specs.PropertyReference{
											Resource: template.InputResource,
											Path:     "nested.first",
										},
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
								"second": {
									Position: 2,
									Name:     "second",
									Path:     "nested.second",
									Label:    labels.Optional,
									Template: specs.Template{
										Reference: &specs.PropertyReference{
											Resource: template.InputResource,
											Path:     "nested.second",
										},
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
							},
						},
					},
					"repeating_string": {
						Position: 9,
						Name:     "repeating_string",
						Path:     "repeating_string",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating_string",
							},
							Repeated: specs.Repeated{
								{
									Scalar: &specs.Scalar{
										Type: types.String,
									},
								},
							},
						},
					},
					"repeating_enum": {
						Position: 10,
						Name:     "repeating_enum",
						Path:     "repeating_enum",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating_enum",
							},
							Repeated: specs.Repeated{
								{
									Enum: enum,
								},
							},
						},
					},
					"repeating_numeric": {
						Position: 11,
						Name:     "repeating_numeric",
						Path:     "repeating_numeric",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating_numeric",
							},
							Scalar: &specs.Scalar{
								Type: types.Int32,
							},
						},
					},
					"repeating": {
						Position: 12,
						Name:     "repeating",
						Path:     "repeating",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating",
							},
							Repeated: specs.Repeated{
								{
									Message: specs.Message{
										"value": {
											Position: 1,
											Name:     "value",
											Path:     "repeating.value",
											Label:    labels.Optional,
											Template: specs.Template{
												Reference: &specs.PropertyReference{
													Resource: template.InputResource,
													Path:     "repeating.value",
												},
												Scalar: &specs.Scalar{
													Type: types.String,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	enum = &specs.Enum{
		Keys: map[string]*specs.EnumValue{
			"UNKNOWN": {
				Key:      "UNKNOWN",
				Position: 0,
			},
			"PENDING": {
				Key:      "PENDING",
				Position: 1,
			},
		},
		Positions: map[int32]*specs.EnumValue{
			0: {
				Key:      "UNKNOWN",
				Position: 0,
			},
			1: {
				Key:      "PENDING",
				Position: 1,
			},
		},
	}
)
