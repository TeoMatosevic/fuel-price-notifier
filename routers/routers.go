package routers

import (
    "fuel-price-notifier/context"
    "net/http"
    "strings"
    "time"
    "crypto/tls"
)

func Auth(ctx context.Context, next func(context.Context) http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        reqToken := r.Header.Get("Authorization")
        if reqToken == "" {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        token := strings.Split(reqToken, "Token ")[1]

        user, err := ctx.Users().FindToken(token)

        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        var tokenEnd int64 = user.Token.End

        if tokenEnd < time.Now().Unix() {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Unauthorized"))
            return
        }

        next(ctx)(w, r)
    }
}

func EnableCors(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if r.Method == "OPTIONS" {
            return
        }
        next(w, r)
    }
}


func Init(ctx *context.Context) { 

    cert, err := tls.LoadX509KeyPair("ca.crt", "ca.key")

    if err != nil {
        panic(err)
    }

    config := &tls.Config{Certificates: []tls.Certificate{cert}}

    router := http.NewServeMux()

    initializeUsersRouters(ctx, router)
    initializeGasStationsRouters(ctx, router)

    server := &http.Server{
        Addr:      ":8080",
        Handler:   router,
        TLSConfig: config,
    }

    err = server.ListenAndServeTLS("", "")

    if err != nil {
        panic(err)
    }
}
