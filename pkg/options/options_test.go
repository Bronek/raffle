package options

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

type dummy_entropy_source struct {
	t    *testing.T
	fail int
}

const dummy_entropy = (1 << 32) + 1

func (t *dummy_entropy_source) Read(b []byte) (int, error) {
	if t.fail > 0 {
		t.fail -= 1
		return 0, fmt.Errorf("Dummy error")
	}
	if len(b) < 4 {
		t.t.FailNow()
	}
	b[3] = 1 // That's `dummy_entropy` in 2 passes
	return 4, nil
}

func TestParseSuccess(t *testing.T) {
	var tests = []struct {
		args []string
		conf Config
	}{
		{[]string{},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{""},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"foo"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 10, args: []string{"foo"}}},
		{[]string{"--seed=0"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"-seed=0"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--seed=1"},
			Config{Seed: 1, SeedFromEntropy: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"-seed=1"},
			Config{Seed: 1, SeedFromEntropy: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--seed=-5026403773975906525"},
			Config{Seed: -5026403773975906525, SeedFromEntropy: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--seed=4596354214482412099"},
			Config{Seed: 4596354214482412099, SeedFromEntropy: false, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--multiplier=1"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 10, args: []string{}}},
		{[]string{"--multiplier=100"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 100, N: 10, args: []string{}}},
		{[]string{"-multiplier=100"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 100, N: 10, args: []string{}}},
		{[]string{"--multiplier=10000"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 10000, N: 10, args: []string{}}},
		{[]string{"--N=1"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 1, args: []string{}}},
		{[]string{"--N=20"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 20, args: []string{}}},
		{[]string{"-N=20"},
			Config{Seed: dummy_entropy, SeedFromEntropy: true, Multiplier: 1, N: 20, args: []string{}}},
		{[]string{"--N=20", "--seed=1", "--multiplier=100", "baz"},
			Config{Seed: 1, SeedFromEntropy: false, Multiplier: 100, N: 20, args: []string{"baz"}}},
		{[]string{"-N=20", "-seed=1", "-multiplier=100"},
			Config{Seed: 1, SeedFromEntropy: false, Multiplier: 100, N: 20, args: []string{}}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := Parse("dummy", tt.args, &dummy_entropy_source{t, 1})
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			if output != "" {
				t.Errorf("%q", output)
				return
			}
			if !reflect.DeepEqual(*conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v, got %T, want %T", *conf, tt.conf, conf.args, tt.conf.args)
			}
		})
	}
}

// Various misbehaving entropy sources, for the error handling test
type dummy_entropy_0_one_pass struct{}

func (dummy_entropy_0_one_pass) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	return len(b), nil
}

type dummy_entropy_0_two_passes struct {
	t *testing.T
}

func (t dummy_entropy_0_two_passes) Read(b []byte) (int, error) {
	if len(b) < 4 {
		t.t.FailNow()
	}
	for i := 0; i < 4; i++ {
		b[i] = 0
	}
	return 4, nil
}

type dummy_entropy_0_four_passes struct {
	t *testing.T
}

func (t dummy_entropy_0_four_passes) Read(b []byte) (int, error) {
	if len(b) < 2 {
		t.t.FailNow()
	}
	b[0] = 0
	b[1] = 0
	return 2, nil
}

type dummy_entropy_one_byte_only struct {
	t *testing.T
}

func (t dummy_entropy_one_byte_only) Read(b []byte) (int, error) {
	if len(b) < 1 {
		t.t.FailNow()
	}
	b[0] = 1
	return 1, nil
}

func TestParseErrors(t *testing.T) {
	var tests = []struct {
		args    []string
		entropy io.Reader
		check   func(string, error, *testing.T)
	}{
		{[]string{"-dummy"},
			&dummy_entropy_source{t, 0},
			func(help string, err error, t *testing.T) {
				if len(help) == 0 {
					t.Errorf("Expected help string, got nothing")
				}
			}},
		{[]string{"-multiplier=0"},
			&dummy_entropy_source{t, 0},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Multiplier out of range") {
					t.Errorf("Expected 'Multiplier out of range' error, got: %v", err)
				}
			}},
		{[]string{"-multiplier=-1"},
			&dummy_entropy_source{t, 0},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Multiplier out of range") {
					t.Errorf("Expected 'Multiplier out of range' error, got: %v", err)
				}
			}},
		{[]string{"-multiplier=10001"},
			&dummy_entropy_source{t, 0},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Multiplier out of range") {
					t.Errorf("Expected 'Multiplier out of range' error, got: %v", err)
				}
			}},
		{[]string{"-seed=0"},
			&dummy_entropy_0_one_pass{},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Entropy source failure, 0 returned") {
					t.Errorf("Expected 'Entropy source failure, 0 returned' error, got: %v", err)
				}
			}},
		{[]string{"-seed=0"},
			&dummy_entropy_0_two_passes{t},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Entropy source failure, 0 returned") {
					t.Errorf("Expected 'Entropy source failure, 0 returned' error, got: %v", err)
				}
			}},
		{[]string{"-seed=0"},
			&dummy_entropy_0_four_passes{t},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Entropy source failure, 0 returned") {
					t.Errorf("Expected 'Entropy source failure, 0 returned' error, got: %v", err)
				}
			}},
		{[]string{"-seed=0"},
			&dummy_entropy_one_byte_only{t},
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Failed to read desired size from entropy source") {
					t.Errorf("Expected 'Failed to read desired size from entropy source' error, got: %v", err)
				}
			}},
		{[]string{"-seed=0"},
			&dummy_entropy_source{t, 3}, // Fails 3 times, 4th returns only 4 bytes
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Failed to read desired size from entropy source") {
					t.Errorf("Expected 'Failed to read desired size from entropy source' error, got: %v", err)
				}
			}},
		{[]string{"-seed=0"},
			&dummy_entropy_source{t, 4}, // Fails 4 times
			func(help string, err error, t *testing.T) {
				if !strings.Contains(err.Error(), "Failed to read entropy source") {
					t.Errorf("Expected 'Failed to read entropy source' error, got: %v", err)
				}
			}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			_, help, err := Parse("dummy", tt.args, tt.entropy)
			if err == nil {
				t.Errorf("Expected error, got success")
				return
			}
			tt.check(help, err, t)
		})
	}
}
