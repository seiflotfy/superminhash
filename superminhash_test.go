package superminhash

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	metro "github.com/dgryski/go-metro"
	minhash "github.com/dgryski/go-minhash"
)

func mhash1(b []byte) uint64 { return metro.Hash64(b, 42) }
func mhash2(b []byte) uint64 { return metro.Hash64(b, 1337) }

func TestSimilarity(t *testing.T) {
	var (
		src     = rand.NewSource(time.Now().UnixNano())
		rand    = rand.New(src)
		length  = rand.Int63n(10000)
		s1      = &Signature{values: make([]float64, length)}
		s2      = &Signature{values: make([]float64, length)}
		modRate = rand.Float64()
		numMods float64
	)

	// modify s2
	for i := range s1.values {
		s1.values[i] = rand.Float64()
		s2.values[i] = s1.values[i]
		if rand.Float64() >= modRate {
			s2.values[i]++
			numMods++
		}
	}

	expected := 1.0 - numMods/float64(length)

	if sim, err := s1.Similarity(s2); err != nil {
		t.Errorf("expected no error, got %v", err)
	} else if int(10000*sim) != int(10000*expected) { // because floats suck like that
		t.Errorf("expected similarity (s1, s2) = %10f, got %10f", expected, sim)
	}
}
func TestSimilarityError(t *testing.T) {
	var (
		s1 = &Signature{values: make([]float64, 10)}
		s2 = &Signature{values: make([]float64, 11)}
	)

	if _, err := s1.Similarity(s2); err == nil {
		t.Errorf("expected error, got mil")
	}
}

func TestComplete(t *testing.T) {
	tests := []struct {
		s1 []string
		s2 []string
	}{
		{
			[]string{"hello", "world", "foo", "baz", "bar", "zomg"},
			[]string{"goodbye", "world", "foo", "qux", "bar", "zomg"},
		},
	}

	for _, tt := range tests {
		s1, _ := NewSignature(10)
		s2, _ := NewSignature(10)
		m1 := minhash.NewMinWise(mhash1, mhash2, 10)
		m2 := minhash.NewMinWise(mhash1, mhash2, 10)

		for _, s := range tt.s1 {
			s1.Push([]byte(s))
			m1.Push([]byte(s))

		}

		for _, s := range tt.s2 {
			s2.Push([]byte(s))
			m2.Push([]byte(s))
		}

		sim1, _ := s1.Similarity(s2)
		t.Log(sim1)
		sim2 := m1.Similarity(m2)
		t.Log(sim2)

		fmt.Println(sim1, sim2)
	}
}

func TestComplete2(t *testing.T) {
	var (
		src     = rand.NewSource(time.Now().UnixNano())
		rand    = rand.New(src)
		length  = rand.Int63n(10000)
		s1, _   = NewSignature(256)
		s2, _   = NewSignature(256)
		m1      = minhash.NewMinWise(mhash1, mhash2, 256)
		m2      = minhash.NewMinWise(mhash1, mhash2, 256)
		modRate = rand.Float64()
		numMods = int64(0)
	)

	for i := 0; i < int(length); i++ {
		s := strconv.Itoa(i)
		s1.Push([]byte(s))
		m1.Push([]byte(s))
		if rand.Float64() > modRate {
			s += "_"
			numMods++
		}
		s2.Push([]byte(s))
		m2.Push([]byte(s))
	}

	sim1, _ := s1.Similarity(s2)
	t.Log(sim1)
	sim2 := m1.Similarity(m2)
	t.Log(sim2)

	t.Log(numMods, length, 1-float64(numMods)/float64(length))
}
