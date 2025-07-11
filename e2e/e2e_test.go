package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/k1LoW/runn"
)

func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests in short mode")
	}

	tests := []struct {
		name string
		file string
	}{
		{
			name: "healthz",
			file: "healthz.yml",
		},
		{
			name: "auth",
			file: "auth.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			p := filepath.Join(".", tt.file)
			
			opts := []runn.Option{
				runn.T(t),
				runn.Book(p),
			}

			o, err := runn.New(opts...)
			if err != nil {
				t.Fatal(err)
			}

			if err := o.Run(ctx); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Make sure we're running in the e2e directory
	if _, err := os.Stat("healthz.yml"); os.IsNotExist(err) {
		os.Chdir("e2e")
	}
	
	os.Exit(m.Run())
}