package flow

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

type MockCodec struct{}

func (codec *MockCodec) Marshal(refs.Store) (io.Reader, error) {
	return nil, nil
}

func (codec *MockCodec) Unmarshal(io.Reader, refs.Store) error {
	return nil
}

type caller struct {
	name       string
	Counter    int
	mutex      sync.Mutex
	Err        error
	references []*specs.Property
}

func (caller *caller) References() []*specs.Property {
	return caller.references
}

func (caller *caller) Do(context.Context, refs.Store) error {
	caller.mutex.Lock()
	caller.Counter++
	caller.mutex.Unlock()
	return caller.Err
}

func NewMockFlowManager(caller Call, revert Call) ([]*Node, *Manager) {
	ctx := instance.NewContext()

	nodes := []*Node{
		NewMockNode("first", caller, revert),
		NewMockNode("second", caller, revert),
		NewMockNode("third", caller, revert),
		NewMockNode("fourth", caller, revert),
	}

	nodes[0].Next = []*Node{nodes[1], nodes[2]}

	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[3]}
	nodes[2].Previous = []*Node{nodes[0]}
	nodes[2].Next = []*Node{nodes[3]}

	nodes[3].Previous = []*Node{nodes[1], nodes[2]}

	return nodes, &Manager{
		ctx:        ctx,
		Starting:   []*Node{nodes[0]},
		References: 0,
		Nodes:      len(nodes),
		Ends:       1,
	}
}

func TestNewManager(t *testing.T) {
	tests := map[string][]*Node{
		"default": {
			{
				Name: "first",
			},
			{
				Name: "second",
			},
		},
	}

	for name, nodes := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			manager := NewManager(ctx, name, nodes, nil, nil)
			if manager == nil {
				t.Fatal("unexpected result, expected a manager to be returned")
			}
		})
	}
}

func TestCallFlowManager(t *testing.T) {
	caller := &caller{}
	nodes, manager := NewMockFlowManager(caller, nil)
	err := manager.Do(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	if caller.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, len(nodes))
	}
}

func TestFailFlowManager(t *testing.T) {
	expected := errors.New("something went wrong")
	reverts := 2
	calls := 2

	rollback := &caller{}
	call := &caller{}

	nodes, manager := NewMockFlowManager(call, rollback)

	nodes[2].Call = &caller{Err: expected}

	err := manager.Do(context.Background(), nil)
	if !errors.Is(err, expected) {
		t.Errorf("unexpected result %s, expected %s", err, expected)
	}

	manager.Wait()

	if call.Counter != calls {
		t.Errorf("unexpected counter total %d, expected %d", call.Counter, calls)
	}

	if rollback.Counter != reverts {
		t.Errorf("unexpected rollback counter total %d, expected %d", rollback.Counter, reverts)
	}
}

func TestBeforeDoFlow(t *testing.T) {
	counter := 0
	call := &caller{}
	_, manager := NewMockFlowManager(call, nil)

	manager.BeforeDo = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	err := manager.Do(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected before do function to be called", counter)
	}
}

func TestBeforeDoFlowErr(t *testing.T) {
	expected := errors.New("unexpected error")
	counter := 0
	call := &caller{}
	_, manager := NewMockFlowManager(call, nil)

	manager.BeforeDo = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	err := manager.Do(context.Background(), nil)
	if !errors.Is(err, expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected before do function to be called", counter)
	}
}

func TestAfterDoFlowErr(t *testing.T) {
	expected := errors.New("unexpected error")
	counter := 0
	call := &caller{}
	_, manager := NewMockFlowManager(call, nil)

	manager.AfterDo = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	err := manager.Do(context.Background(), nil)
	if !errors.Is(err, expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", err, expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected before do function to be called", counter)
	}
}

func TestAfterDoFlow(t *testing.T) {
	counter := 0
	call := &caller{}
	_, manager := NewMockFlowManager(call, nil)

	manager.AfterDo = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	err := manager.Do(context.Background(), nil)
	if err != nil {
		t.Error(err)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after do function to be called", counter)
	}
}

func TestBeforeRollbackFlow(t *testing.T) {
	counter := 0
	call := &caller{}
	nodes, manager := NewMockFlowManager(call, nil)

	manager.BeforeRollback = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	manager.wg.Add(1)
	manager.Revert(NewTracker("", len(nodes)), nil)

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected before rollback function to be called", counter)
	}
}

func TestBeforeRollbackFlowErr(t *testing.T) {
	expected := errors.New("unexpected error")
	counter := 0
	call := &caller{}
	nodes, manager := NewMockFlowManager(call, nil)

	manager.BeforeRollback = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	manager.wg.Add(1)
	manager.Revert(NewTracker("", len(nodes)), nil)

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected before rollback function to be called", counter)
	}
}

func TestAfterRollbackFlow(t *testing.T) {
	counter := 0
	call := &caller{}
	nodes, manager := NewMockFlowManager(call, nil)

	manager.AfterRollback = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	manager.wg.Add(1)
	manager.Revert(NewTracker("", len(nodes)), nil)

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after rollback function to be called", counter)
	}
}

func TestAfterRollbackFlowErr(t *testing.T) {
	expected := errors.New("unexpected error")
	counter := 0
	call := &caller{}
	nodes, manager := NewMockFlowManager(call, nil)

	manager.AfterRollback = func(ctx context.Context, manager *Manager, store refs.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	manager.wg.Add(1)
	manager.Revert(NewTracker("", len(nodes)), nil)

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected before rollback function to be called", counter)
	}
}

func TestAfterManagerFunctions(t *testing.T) {
	type test struct {
		expected int
		stack    functions.Stack
	}

	current := 0
	counter := func(store refs.Store) error {
		current++
		return nil
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			stack: functions.Stack{
				"first": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        counter,
					Returns:   &specs.Property{},
				},
			},
		},
		"multiple": {
			expected: 3,
			stack: functions.Stack{
				"first": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        counter,
					Returns:   &specs.Property{},
				},
				"second": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        counter,
					Returns:   &specs.Property{},
				},
				"third": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        counter,
					Returns:   &specs.Property{},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			current = 0
			ctx := instance.NewContext()
			manager := NewManager(ctx, name, []*Node{}, test.stack, nil)
			if manager == nil {
				t.Fatal("unexpected result, expected a manager to be returned")
			}

			store := refs.NewReferenceStore(1)
			err := manager.Do(context.Background(), store)
			if err != nil {
				t.Fatalf("unexpected error, %s", err)
			}

			if current != test.expected {
				t.Fatalf("unexpected count value %d, expected %d", current, test.expected)
			}
		})
	}
}

func TestAfterManagerFunctionsError(t *testing.T) {
	type test struct {
		expected int
		stack    functions.Stack
	}

	expected := errors.New("unexpected error")

	pass := func(store refs.Store) error {
		return nil
	}

	fail := func(store refs.Store) error {
		return expected
	}

	tests := map[string]test{
		"single": {
			expected: 1,
			stack: functions.Stack{
				"first": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        fail,
					Returns:   &specs.Property{},
				},
			},
		},
		"multiple": {
			expected: 3,
			stack: functions.Stack{
				"first": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        pass,
					Returns:   &specs.Property{},
				},
				"second": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        pass,
					Returns:   &specs.Property{},
				},
				"third": &functions.Function{
					Arguments: []*specs.Property{},
					Fn:        fail,
					Returns:   &specs.Property{},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			manager := NewManager(ctx, name, []*Node{}, test.stack, nil)
			if manager == nil {
				t.Fatal("unexpected result, expected a manager to be returned")
			}

			store := refs.NewReferenceStore(1)
			err := manager.Do(context.Background(), store)
			if err == nil {
				t.Fatal("unexpected pass expected a error to be returned")
			}

			if !errors.Is(err, expected) {
				t.Fatalf("unexpected err '%s', expected the expected error to be returned '%s'", err, expected)
			}
		})
	}
}
