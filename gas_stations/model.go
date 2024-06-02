package gas_stations

import (
    "errors"
    "fuel-price-notifier/location"
    "strings"
    "time"
    "math"
)

const (
    Gasoline = iota
    Diesel = iota
    LPG = iota
)

type InsertedGasStation struct {
    Name string 
    Address string
    Location *location.Location
    FuelPrices map[string]float64
}

type GasStation struct {
    name string
    location *location.Location
    address string
    fuelPrices map[int]float64
    fuelPricesHistory map[int64]map[int]float64
}

func (g *GasStation) Location() *location.Location {
    return g.location
}

func (g *GasStation) FuelPrices() map[int]float64 {
    return g.fuelPrices
}

func (g *GasStation) FuelPricesHistory() map[int64]map[int]float64 {
    return g.fuelPricesHistory
}

func (g *GasStation) ToDto() *GasStationDto {
    fuelPrices := make(map[string]float64)
    for k, v := range g.fuelPrices {
        if k == Gasoline {
            fuelPrices["gasoline"] = RoundF(v, 3)
        }
        if k == Diesel {
            fuelPrices["diesel"] = RoundF(v, 3)
        }
        if k == LPG {
            fuelPrices["lpg"] = RoundF(v, 3)
        }
    }
    return &GasStationDto{g.name, g.address, g.location, fuelPrices}
}

func Create(i *InsertedGasStation) (*GasStation, error) {
    l := location.Create(i.Location.Latitude, i.Location.Longitude)

    if !l.Verify() {
        return nil, errors.New("Invalid location")
    }

    fuelPrices := make(map[int]float64)
    fuelPricesCopy := make(map[int]float64)
    fuelPricesHistory := make(map[int64]map[int]float64)
    for k, v := range i.FuelPrices {
        if strings.ToLower(k) == "gasoline" {
            if _, keyExists := fuelPrices[Gasoline]; keyExists {
                return nil, errors.New("Duplicate fuel type")
            }
            fuelPrices[Gasoline] = v
            fuelPricesCopy[Gasoline] = v
        } else if strings.ToLower(k) == "diesel" {
            if _, keyExists := fuelPrices[Diesel]; keyExists {
                return nil, errors.New("Duplicate fuel type")
            }
            fuelPrices[Diesel] = v
            fuelPricesCopy[Diesel] = v
        } else if strings.ToLower(k) == "lpg" {
            if _, keyExists := fuelPrices[LPG]; keyExists {
                return nil, errors.New("Duplicate fuel type")
            }
            fuelPrices[LPG] = v
            fuelPricesCopy[LPG] = v
        } else {
            return nil, errors.New("Invalid fuel type")
        }
    }

    current := time.Now().Unix()
    fuelPricesHistory[current] = fuelPricesCopy

    return &GasStation{i.Name, l, i.Address, fuelPrices, fuelPricesHistory}, nil
}

type GasStationDto struct {
    Name string `json:"name"`
    Address string `json:"address"`
    Location *location.Location `json:"location"`
    FuelPrices map[string]float64 `json:"fuel_prices"`
}

type GasStations struct {
    G map[string]*GasStation
}

func (g *GasStations) UpdateFuelPrices(a string, fuelPrices map[int]float64, t int64) {
    fuelPricesCopy := make(map[int]float64)
    for k, v := range fuelPrices {
        fuelPricesCopy[k] = v
    }
    g.G[a].fuelPrices = fuelPrices
    g.G[a].fuelPricesHistory[t] = fuelPricesCopy
}

func (g *GasStations) Add(gs *GasStation) error {
    if _, keyExists := g.G[gs.address]; keyExists {
        return errors.New("Gas station already exists")
    }
    g.G[gs.address] = gs
    return nil
}

func (g *GasStations) ToDto() []*GasStationDto {
    gasStations := make([]*GasStationDto, 0, len(g.G))
    for _, gasStation := range g.G {
        fuelPrices := make(map[string]float64)
        for k, v := range gasStation.fuelPrices {
            if k == Gasoline {
                fuelPrices["gasoline"] = RoundF(v, 3)
            } else if k == Diesel {
                fuelPrices["diesel"] = RoundF(v, 3)
            } else if k == LPG {
                fuelPrices["lpg"] = RoundF(v, 3)
            }
        }
        gasStations = append(gasStations, &GasStationDto{gasStation.name, gasStation.address, gasStation.location, fuelPrices})
    }
    return gasStations
}

func (g *GasStations) Get(address string) (*GasStation, error) {
    gasStation, keyExists := g.G[address]
    if !keyExists {
        return nil, errors.New("Gas station not found")
    }
    return gasStation, nil
}

func RoundF(v float64, p uint) float64 {
    r := math.Pow(10, float64(p))
    return math.Round(v * r) / r
}

func GetFuelType(i int) string {
    if i == Gasoline {
        return "gasoline"
    } else if i == Diesel {
        return "diesel"
    } else if i == LPG {
        return "lpg"
    }
    panic("Invalid fuel type")
}

