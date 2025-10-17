package data

import (
	"testing"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/validator"
)

func TestValidateInstance(t *testing.T) {
	cases := []struct {
		name     string
		instance Instance
		expectOK bool
		expect   map[string]string
	}{
		{
			name: "valid_min_boundary",
			instance: Instance{
				Name:        "srv-a",
				Ip:          "1.2.3.4",
				RefreshRate: 1 * time.Minute,
			},
			expectOK: true,
			expect:   map[string]string{},
		},
		{
			name: "valid_max_boundary",
			instance: Instance{
				Name:        "srv-b",
				Ip:          "5.6.7.8",
				RefreshRate: 24 * time.Hour,
			},
			expectOK: true,
			expect:   map[string]string{},
		},
		{
			name: "empty_name",
			instance: Instance{
				Name:        "",
				Ip:          "1.1.1.1",
				RefreshRate: 5 * time.Minute,
			},
			expectOK: false,
			expect: map[string]string{
				"name": "must be provided",
			},
		},
		{
			name: "name_too_long",
			instance: Instance{
				Name:        makeString(501),
				Ip:          "1.1.1.1",
				RefreshRate: 5 * time.Minute,
			},
			expectOK: false,
			expect: map[string]string{
				"name": "name must not be longer than 500 bytes",
			},
		},
		{
			name: "empty_ip",
			instance: Instance{
				Name:        "srv",
				Ip:          "",
				RefreshRate: 5 * time.Minute,
			},
			expectOK: false,
			expect: map[string]string{
				"ip": "must be provided",
			},
		},
		{
			name: "refresh_too_small",
			instance: Instance{
				Name:        "srv",
				Ip:          "1.1.1.1",
				RefreshRate: 30 * time.Second,
			},
			expectOK: false,
			expect: map[string]string{
				"refresh_rate": "refresh rate must be at least one minute",
			},
		},
		{
			name: "refresh_too_large",
			instance: Instance{
				Name:        "srv",
				Ip:          "1.1.1.1",
				RefreshRate: 25 * time.Hour,
			},
			expectOK: false,
			expect: map[string]string{
				"refresh_rate": "refresh rate must be less than 24 hours",
			},
		},
		{
			name: "multiple_errors",
			instance: Instance{
				Name:        "",
				Ip:          "",
				RefreshRate: 0,
			},
			expectOK: false,
			expect: map[string]string{
				"name":         "must be provided",
				"ip":           "must be provided",
				"refresh_rate": "refresh rate must be at least one minute",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			inst := tc.instance // copy
			ValidateInstance(v, &inst)

			if v.Valid() != tc.expectOK {
				t.Fatalf("Valid()=%v, want %v with errors=%v", v.Valid(), tc.expectOK, v.Errors)
			}

			if !tc.expectOK {
				if len(v.Errors) != len(tc.expect) {
					t.Fatalf("len(errors)=%d, want %d; got=%v", len(v.Errors), len(tc.expect), v.Errors)
				}
				for key, want := range tc.expect {
					got, ok := v.Errors[key]
					if !ok {
						t.Fatalf("expected error for key %q not found; got=%v", key, v.Errors)
					}
					if got != want {
						t.Fatalf("error for key %q=%q, want %q", key, got, want)
					}
				}
			}
		})
	}
}

func makeString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
