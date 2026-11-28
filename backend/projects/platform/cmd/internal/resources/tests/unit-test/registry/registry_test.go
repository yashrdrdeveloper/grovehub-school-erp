package registry_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"grovehub/cmd/internal/src/command"
	"grovehub/cmd/internal/src/registry"
)

type fakeCommand struct {
	args []string
}

func (f *fakeCommand) Run(args []string, output, diagnostics io.Writer) int {
	f.args = args
	_, _ = io.WriteString(output, "output")
	_, _ = io.WriteString(diagnostics, "diagnostic")
	return 7
}

func TestDispatch(t *testing.T) {
	fake := &fakeCommand{}
	r := registry.New(map[string]command.Command{"serve": fake, "migrate": fake, "nil": nil})
	var output, diagnostics bytes.Buffer
	if code := r.Run(nil, &output, &diagnostics); code != 7 || len(fake.args) != 0 {
		t.Fatalf("default command: %d, %v", code, fake.args)
	}
	if code := r.Run([]string{"migrate", "all"}, &output, &diagnostics); code != 7 || !reflect.DeepEqual(fake.args, []string{"all"}) {
		t.Fatalf("explicit command: %d, %v", code, fake.args)
	}
	for _, name := range []string{"unknown", "nil"} {
		diagnostics.Reset()
		if code := r.Run([]string{name}, &output, &diagnostics); code != 2 || diagnostics.Len() == 0 {
			t.Fatalf("invalid command %q: %d", name, code)
		}
	}
}
