package alias_sample

import (
	"log"
	"math"
	"testing"

	"pgregory.net/rapid"
)

func TestInit(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		probs := rapid.SliceOfN(rapid.Float64Range(0.001, 5.0), 1, 100).Draw(t, "probs")
		as, err := Init(probs)
		if err != nil {
			log.Fatalf("got err %v\n", err)
		}
		sz := 1_000_000
		res := make([]int, len(probs))
		for range sz {
			r := as.Next()
			res[r] += 1
		}

		norm_probs := make([]float64, len(probs))
		copy(norm_probs, probs)

		var tot float64
		for _, p := range norm_probs {
			tot += p
		}

		for i := range norm_probs {
			norm_probs[i] /= tot
		}

		res_probs := make([]float64, len(probs))
		for i := range len(probs) {
			res_probs[i] = float64(res[i]) / float64(sz)
		}
		for i, p := range res_probs {
			if math.Abs(p-norm_probs[i]) > 0.01 {
				t.Fatalf("failed: %f, %f, %v %v %v\n", p, norm_probs[i], norm_probs, res, res_probs)
			}
		}
	})
}
