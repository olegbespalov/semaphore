package compare

import (
	"testing"

	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func TestCompareMessages(t *testing.T) {
	createScalar := func() *specs.Property {
		return &specs.Property{
			Name:     "age",
			Path:     "dog",
			Position: 0,
			Label:    labels.Required,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.Int32,
				},
			},
		}
	}
	messageA := specs.Message{
		"age": createScalar(),
	}
	messageB := specs.Message{
		"age":   createScalar(),
		"color": createScalar(),
	}

	type args struct {
		expected specs.Message
		given    specs.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"should not match different messages",
			args{messageA, messageB},
			true,
		},

		{
			"should not match messages against nil",
			args{nil, messageA},
			true,
		},

		{
			"should not match nil against messages",
			args{messageA, nil},
			true,
		},

		{
			"should match",
			args{messageA, messageA},
			false,
		},

		{
			"should match nils",
			args{nil, nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompareMessages(tt.args.expected, tt.args.given); (err != nil) != tt.wantErr {
				t.Errorf("CompareMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
