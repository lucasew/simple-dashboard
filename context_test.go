package godashboard

import (
	"errors"
	"testing"
)

func TestFirstCPUPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		vals    []float64
		err     error
		want    float64
		wantErr bool
	}{
		{name: "propagates error", err: errors.New("boom"), wantErr: true},
		{name: "empty slice", vals: nil, wantErr: true},
		{name: "zero-length", vals: []float64{}, wantErr: true},
		{name: "single sample", vals: []float64{12.5}, want: 12.5},
		{name: "uses first of many", vals: []float64{1, 2, 3}, want: 1},
		// error wins even if vals are present
		{name: "error with vals", vals: []float64{9}, err: errors.New("nope"), wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := firstCPUPercent(tt.vals, tt.err)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("firstCPUPercent() err = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("firstCPUPercent() unexpected err: %v", err)
			}
			if got != tt.want {
				t.Fatalf("firstCPUPercent() = %v, want %v", got, tt.want)
			}
		})
	}
}
