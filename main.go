package main

import (
    "fuel-price-notifier/context"
    "fuel-price-notifier/routers"
    "fuel-price-notifier/fluctuations"
)

func main() {
    var ctx context.Context = context.NewNonPersistantContext()
    fluctuations.Init(&ctx)
    routers.Init(&ctx)
}
