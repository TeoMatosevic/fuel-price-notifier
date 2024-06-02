package fluctuations

import (
	"fuel-price-notifier/context"
	"time"
)

const FLUCTUATION_INTERVAL = 20

type Fluctuation interface {
    Fluctuate(n float64) float64
}

func notifier(ctx context.Context, ch chan map[string]map[int]float64, f Fluctuation) {
    for {
        gs := ctx.GasStations()
        fuelPricesMap := make(map[string]map[int]float64)
        for k, v := range gs.G {
            fp := v.FuelPrices()
            for k, v := range fp {
                fp[k] = f.Fluctuate(v)
            }
            fuelPricesMap[k] = fp
        }
        ch <- fuelPricesMap
        time.Sleep(FLUCTUATION_INTERVAL * time.Second)
    }
}

func updater(ctx context.Context, ch chan map[string]map[int]float64) {
    for {
        fuelPrices := <-ch
        gasStations := ctx.GasStations()
        curr := time.Now().Unix()
        for k, v := range fuelPrices {
            ctx.LookGasStationsMutex()
            gasStations.UpdateFuelPrices(k, v, curr)
            ctx.UnlockGasStationsMutex()
        }
    }
}

func Init(ctx *context.Context) {
    ch := make(chan map[string]map[int]float64)
    f := NormalDistributionFluctuation{}
    go notifier(*ctx, ch, &f)
    go updater(*ctx, ch)
}

