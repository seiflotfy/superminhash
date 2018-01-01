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
	p      []uint16
	q      []int64
	i      int64
}

// NewSignature ...
func NewSignature(length uint16) (*Signature, error) {
	if length == 0 {
		return nil, fmt.Errorf("length has to be >= 1")
	}
	values := make([]float32, length)
	p := make([]uint16, length)
	q := make([]int64, length)
	for i := range values {
		values[i] = math.MaxUint32
		q[i] = -1
		p[i] = uint16(i)
	}
	return &Signature{
		values: values,
		p:      p,
		q:      q,
		i:      0,
	}, nil
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

	for j := range sig.values {
		r := float32(rnd.Next()) / float32(math.MaxUint32)
		offset := rnd.Next() % uint32(len(sig.values)-j)
		k := uint32(j) + offset

		if sig.q[j] != sig.i {
			sig.q[j] = sig.i
			sig.p[j] = uint16(j)
		}

		if sig.q[k] != sig.i {
			sig.q[k] = sig.i
			sig.p[k] = uint16(k)
		}

		sig.p[j], sig.p[k] = sig.p[k], sig.p[j]
		newVal := r + float32(sig.p[j])
		if newVal < sig.values[j] {
			sig.values[j] = newVal
		}
	}
	sig.i++
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
