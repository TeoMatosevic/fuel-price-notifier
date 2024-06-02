package routers

import (
    "net/http"
    "fuel-price-notifier/context"
    "fuel-price-notifier/users"
    "encoding/json"
)

func allUsers(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            ctx.LookUsersMutex()
            defer ctx.UnlockUsersMutex()

            users := ctx.Users().ToDto()

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(users)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func getUser(w http.ResponseWriter, _ *http.Request, user *users.User) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user.ToDto())
}

func deleteUser(w http.ResponseWriter, _ *http.Request, user *users.User, u *users.Users) {
    err := u.Delete(user)

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.WriteHeader(http.StatusOK)
}

func modifyUser(w http.ResponseWriter, r *http.Request, user *users.User) {
    name := r.PostFormValue("name")
    password := r.PostFormValue("password")

    if name == "" {
        name = user.Name
    }

    var passwordHash []byte

    if password == "" {
        passwordHash = user.PasswordHash
    } else {
        passwordHash = users.Hash(password)
    }

    user.Modify(name, passwordHash)

    w.WriteHeader(http.StatusOK)
}


func user(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        email := r.PathValue("email")

        ctx.LookUsersMutex()
        defer ctx.UnlockUsersMutex()

        user, err := ctx.Users().Get(email)

        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(err.Error()))
            return
        }
        switch r.Method {
        case http.MethodGet:
            getUser(w, r, user)
        case http.MethodDelete:
            deleteUser(w, r, user, ctx.Users())
        case http.MethodPut:
            modifyUser(w, r, user)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func login(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            email := r.PostFormValue("email")
            password := r.PostFormValue("password")

            if email == "" || password == "" {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte("Invalid input"))
                return
            }

            pHash := users.Hash(password)

            ctx.LookUsersMutex()
            defer ctx.UnlockUsersMutex()

            user, err := ctx.Users().Get(email)

            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(err.Error()))
                return
            }

            if string(pHash) != string(user.PasswordHash) {
                w.WriteHeader(http.StatusUnauthorized)
                w.Write([]byte("Invalid credentials"))
                return
            }

            token := users.CreateToken()

            user.SetToken(token)

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(map[string]string{"token": token.Verification})
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func register(ctx context.Context) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            name := r.PostFormValue("name")
            email := r.PostFormValue("email")
            password := r.PostFormValue("password")

            if name == "" || email == "" || password == "" {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte("Invalid input"))
                return
            }

            pHash := users.Hash(password)

            user := users.Create(name, email, pHash)

            ctx.LookUsersMutex()
            defer ctx.UnlockUsersMutex()
            
            err := ctx.Users().Add(user)

            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte(err.Error()))
                return
            }
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    }
}

func initializeUsersRouters(ctx *context.Context, h *http.ServeMux) {
    h.HandleFunc("/api/v1/login", EnableCors(login(*ctx)))
    h.HandleFunc("/api/v1/register", EnableCors(register(*ctx)))
    h.HandleFunc("/api/v1/users/{email}", EnableCors(Auth(*ctx, user)))
    h.HandleFunc("/api/v1/users", EnableCors(Auth(*ctx, allUsers)))
}
