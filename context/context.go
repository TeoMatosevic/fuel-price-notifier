package context

import (
	"fuel-price-notifier/gas_stations"
	"fuel-price-notifier/users"
)

type Context interface {
    Users() *users.Users
    LookUsersMutex()
    UnlockUsersMutex()
    GasStations() *gas_stations.GasStations
    LookGasStationsMutex()
    UnlockGasStationsMutex()
}
