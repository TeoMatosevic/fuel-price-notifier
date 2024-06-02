package context

import (
	"fuel-price-notifier/gas_stations"
	"fuel-price-notifier/users"
	"sync"
)

type NonPersistantContext struct {
    u *users.Users
    mu *sync.Mutex
    g *gas_stations.GasStations
    mg *sync.Mutex
}

func (c *NonPersistantContext) Users() *users.Users {
    return c.u
}

func (c *NonPersistantContext) LookUsersMutex() {
    c.mu.Lock()
}

func (c *NonPersistantContext) UnlockUsersMutex() {
    c.mu.Unlock()
}

func (c *NonPersistantContext) GasStations() *gas_stations.GasStations {
    return c.g
}

func (c *NonPersistantContext) LookGasStationsMutex() {
    c.mg.Lock()
}

func (c *NonPersistantContext) UnlockGasStationsMutex() {
    c.mg.Unlock()
}

func NewNonPersistantContext() Context {
    mu := sync.Mutex{}
    mg := sync.Mutex{}
    u := users.Users{U: make(map[string]*users.User)}
    g := gas_stations.GasStations{G: make(map[string]*gas_stations.GasStation)}
    return &NonPersistantContext{&u, &mu, &g, &mg}
}
