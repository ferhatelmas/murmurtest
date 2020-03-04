package murmurtest

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	m1 "github.com/spaolacci/murmur3"
	m2 "github.com/twmb/murmur3"
)

const (
	size = 1023
	base = 10000
)

type buf []byte

func (b buf) Generate(r *rand.Rand, _ int) reflect.Value {
	n := rand.Intn(size)
	ls := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		ls = append(ls, byte(rand.Intn(256)))
	}
	return reflect.ValueOf(ls)
}

func TestSum(t *testing.T) {
	rand.Seed(time.Now().Unix())

	tests := []struct {
		name   string
		runner interface{}
		scale  float64
	}{
		{
			name:  "sum32",
			scale: base * 4,
			runner: func(b buf) bool {
				return m1.Sum32(b) == m2.Sum32(b)
			},
		},
		{
			name:  "sum64",
			scale: base * 2,
			runner: func(b buf) bool {
				return m1.Sum64(b) == m2.Sum64(b)
			},
		},
		{
			name:  "sum128",
			scale: base,
			runner: func(b buf) bool {
				h1a, h1b := m1.Sum128(b)
				h2a, h2b := m2.Sum128(b)
				return h1a == h2a && h1b == h2b
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := quick.Check(tt.runner, &quick.Config{MaxCountScale: tt.scale}); err != nil {
				t.Error(err)
			}
		})
	}
}
