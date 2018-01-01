package superminhash

import (
	"fmt"
	"math"

	"github.com/dgryski/go-pcgr"

	"github.com/dgryski/go-metro"
)

// Signature ...
type Signature struct {
	values []float32
}

// NewSignature ...
func NewSignature(length uint8) (*Signature, error) {
	if length == 0 {
		return nil, fmt.Errorf("length has to be >= 1")
	}
	values := make([]float32, length)
	for i := range values {
		values[i] = math.MaxUint32
	}
	return &Signature{values: values}, nil
}

// Length ...
func (sig *Signature) Length() int {
	return len(sig.values)
}

// Push ...
func (sig *Signature) Push(b []byte) {
	x := metro.Hash64(b, 42)
	rnd := pcgr.New(int64(x), 0)
	for i, v := range sig.values {
		r := float32(rnd.Next()) / float32(math.MaxUint32)
		if r < v {
			sig.values[i] = r
		}
	}
}

// Similarity ...
func (sig *Signature) Similarity(other *Signature) (float64, error) {
	sim := 0.0
	if len(sig.values) != len(other.values) {
		return 0, fmt.Errorf("signatures not of same length, sign has length %d, while other has length %d", len(sig.values), len(other.values))
	}
	for i, element := range sig.values {
		// fmt.Println(element, other.values[i], element == other.values[i])
		if element == other.values[i] {
			sim++
		}
	}
	return sim / float64(len(sig.values)), nil
}
