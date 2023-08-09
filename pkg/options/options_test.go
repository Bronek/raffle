package options

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseSuccess(t *testing.T) {
	var tests = []struct {
		args []string
		conf Config
	}{
		{[]string{},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{""},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"foo"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 10, args: []string{"foo"}}},
		{[]string{"--seed=0"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"-seed=0"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--seed=1"},
			Config{Seed: 1, SeedFromSystem: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"-seed=1"},
			Config{Seed: 1, SeedFromSystem: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--seed=-5026403773975906525"},
			Config{Seed: -5026403773975906525, SeedFromSystem: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--seed=4596354214482412099"},
			Config{Seed: 4596354214482412099, SeedFromSystem: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--multiplier=1"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--multiplier=100"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 100, N: 10, args: []string{}}},
		{[]string{"-multiplier=100"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 100, N: 10, args: []string{}}},
		{[]string{"--multiplier=10000"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 10000, N: 10, args: []string{}}},
		{[]string{"--N=1"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 1, args: []string{}}},
		{[]string{"--N=20"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 20, args: []string{}}},
		{[]string{"-N=20"},
			Config{Seed: 0, SeedFromSystem: true, Multiplier: 1, N: 20, args: []string{}}},
		{[]string{"--N=20", "--seed=1", "--multiplier=100", "baz"},
			Config{Seed: 1, SeedFromSystem: false, Multiplier: 100, N: 20, args: []string{"baz"}}},
		{[]string{"-N=20", "-seed=1", "-multiplier=100"},
			Config{Seed: 1, SeedFromSystem: false, Multiplier: 100, N: 20, args: []string{}}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := Parse("dummy", tt.args)
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			if output != "" {
				t.Errorf("%q", output)
				return
			}
			if tt.conf.SeedFromSystem {
				// Expect that conf.Seed was populated by system entropy
				if conf.Seed == 0 {
					t.Errorf("Unexpected random value 0 from system entropy")
					return
				}
				conf.Seed = 0
			}
			if !reflect.DeepEqual(*conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v, got %T, want %T", *conf, tt.conf, conf.args, tt.conf.args)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	var tests = []struct {
		args  []string
		check func(string, error, *testing.T)
	}{
		{[]string{"--dummy"},
			func(help string, err error, t *testing.T) {
				if len(help) == 0 {
					t.Errorf("Expected help string, got nothing")
				}
			}},
		{[]string{"--multiplier=0"},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "out of range") {
					t.Errorf("Expected 'out of range' error")
				}
			}},
		{[]string{"--multiplier=-1"},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "out of range") {
					t.Errorf("Expected 'out of range' error")
				}
			}},
		{[]string{"--multiplier=10001"},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "out of range") {
					t.Errorf("Expected 'out of range' error")
				}
			}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			_, help, err := Parse("dummy", tt.args)
			if err == nil {
				t.Errorf("Expected error, got success")
				return
			}
			tt.check(help, err, t)
		})
	}
}
