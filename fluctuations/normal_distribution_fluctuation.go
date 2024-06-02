package fluctuations

import (
	"math/rand/v2"
)

type NormalDistributionFluctuation struct {}

func (r *NormalDistributionFluctuation) Fluctuate(n float64) float64 {
    rnd := rand.NormFloat64() * 0.5 + 100

    res := n * rnd / 100

    return res
}

