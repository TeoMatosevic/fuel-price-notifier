package routers

import (
	"encoding/json"
	"fmt"
	"fuel-price-notifier/context"
	"fuel-price-notifier/gas_stations"
	"fuel-price-notifier/location"
	"net/http"
	"strconv"
	"time"
)

func gasStations(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx.LookGasStationsMutex()
        defer ctx.UnlockGasStationsMutex()

        switch r.Method {
        case http.MethodGet:
            getGasStations(w, r, ctx.GasStations())
        case http.MethodPost:
            createGasStation(w, r, ctx.GasStations())
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func createGasStation(w http.ResponseWriter, r *http.Request, g *gas_stations.GasStations) { 
    var gs gas_stations.InsertedGasStation

    err := json.NewDecoder(r.Body).Decode(&gs)

    if err != nil {
        fmt.Println(err)
        return
    }
    
    gasStation, err := gas_stations.Create(&gs)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }
    
    err = g.Add(gasStation)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func getGasStations(w http.ResponseWriter, _ *http.Request, g *gas_stations.GasStations) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(g.ToDto())
}

func handleClosestGasStations(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx.LookGasStationsMutex()
        defer ctx.UnlockGasStationsMutex()

        switch r.Method {
        case http.MethodGet:
            getThreeClosestGasStations(w, r, ctx.GasStations())
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func getThreeClosestGasStations(w http.ResponseWriter, r *http.Request, g *gas_stations.GasStations) {
    var l location.Location

    lon, errLon := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
    lat, errLat := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)

    if errLon != nil || errLat != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid location"))
        return
    }

    l.Latitude = lat
    l.Longitude = lon

    if !l.Verify() {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid location"))
        return
    }

    stations := make([]*gas_stations.GasStation, 0)
    
    for _, gasStation := range g.G {
        if len(stations) < 3 {
            stations = append(stations, gasStation)
        } else {
            max_distance_station := stations[0]
            for _, station := range stations {
                if station.Location().DistanceTo(&l) > max_distance_station.Location().DistanceTo(&l) {
                    max_distance_station = station
                }
            }

            if gasStation.Location().DistanceTo(&l) < max_distance_station.Location().DistanceTo(&l) {
                for i, station := range stations {
                    if station == max_distance_station {
                        stations = append(stations[:i], stations[i+1:]...)
                        break
                    }
                }
                stations = append(stations, gasStation)
            }

        }
    }

    result := make([]*gas_stations.GasStationDto, 0)

    for _, gasStation := range stations {
        result = append(result, gasStation.ToDto())
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(result)
}

func fuelPriceHistory(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx.LookGasStationsMutex()
        defer ctx.UnlockGasStationsMutex()

        switch r.Method {
        case http.MethodGet:
            getFuelPriceHistory(w, r, ctx.GasStations())
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func getFuelPriceHistory(w http.ResponseWriter, r *http.Request, g *gas_stations.GasStations) {
    address := r.URL.Query().Get("address")

    if address == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid address"))
        return
    }
    
    gasStation, err := g.Get(address)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    ph := gasStation.FuelPricesHistory()

    res := make(map[string]map[string]float64)

    for k, v := range ph {
        fp := make(map[string]float64)

        for k1, v1 := range v {
            fp[gas_stations.GetFuelType(k1)] = gas_stations.RoundF(v1, 3)
        }

        t := time.Unix(k, 0)
        res[t.Format(time.RFC3339Nano)] = fp
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}


func initializeGasStationsRouters(ctx *context.Context, h *http.ServeMux) {
    h.HandleFunc("/api/v1/gas-stations",EnableCors(Auth(*ctx, gasStations)))
    h.HandleFunc("/api/v1/gas-stations/closest",EnableCors(Auth(*ctx, handleClosestGasStations)))
    h.HandleFunc("/api/v1/gas-stations/price-history",EnableCors(Auth(*ctx, fuelPriceHistory)))
}

