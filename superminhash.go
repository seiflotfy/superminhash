package superminhash

import (
	"fmt"
	"math"

	metro "github.com/dgryski/go-metro"
	pcgr "github.com/dgryski/go-pcgr"
)

// Signature ...
type Signature struct {
	values []float64
	p      []uint16
	q      []int64
	b      []int64
	i      int64
	a      uint16
}

// NewSignature ...
func NewSignature(length uint16) (*Signature, error) {
	if length == 0 {
		return nil, fmt.Errorf("length has to be >= 1")
	}
	values := make([]float64, length)
	p := make([]uint16, length)
	q := make([]int64, length)
	b := make([]int64, length)
	for i := range values {
		values[i] = math.MaxUint32
		q[i] = -1
		p[i] = uint16(i)
		b[i] = 0
	}
	b[length-1] = int64(length)
	return &Signature{
		values: values,
		p:      p,
		q:      q,
		b:      b,
		i:      0,
		a:      uint16(length) - 1,
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

	for j := uint16(0); j < sig.a; j++ {
		r := float64(rnd.Next()) / float64(math.MaxUint32)
		offset := rnd.Next() % uint32(uint16(len(sig.values))-j)
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
		rj := r + float64(j)
		if rj < sig.values[sig.p[j]] {
			jc := uint16(math.Min(sig.values[sig.p[j]], float64(len(sig.values)-1)))
			sig.values[sig.p[j]] = rj
			if j < jc {
				sig.b[jc]--
				sig.b[j]++
				for sig.b[sig.a] == 0 {
					sig.a--
				}
			}
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
		if element == other.values[i] {
			sim++
		}
	}
	return sim / float64(len(sig.values)), nil
}
