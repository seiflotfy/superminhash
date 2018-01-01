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
	// initialize pseudo-random generator with seed d
	d := metro.Hash64(b, 42)
	rnd := pcgr.New(int64(d), 0)

	// (p0,p1,...,pm−1)←(0,1,...,m−1)
	p := make([]uint8, len(sig.values))
	for i := range p {
		p[i] = uint8(i)
	}

	for j := range p {
		offset := rnd.Next() % uint32(len(p)-j)
		k := uint32(j) + offset
		p[j], p[k] = p[k], p[j]
	}

	for j, v := range sig.values {
		r := float32(rnd.Next()) / float32(math.MaxUint32)
		newVal := r + float32(p[j])
		if newVal < v {
			sig.values[j] = newVal
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
