package users

import (
    "errors"
    "time"
)

const duration = 60 * 60 * 24

type Token struct {
    Start int64
    End int64 
    Verification string 
}

func CreateToken() *Token {
    start := time.Now().Unix()
    end := start + duration
    v := GenerateRandomString(32)
    return &Token{start, end, v}
}

type User struct {
    Name string
    Email string
    PasswordHash []byte
    Token *Token
}

func (u *User) SetToken(t *Token) {
    u.Token = t
}

func (u *User) Modify(name string, password []byte) {
    u.Name = name
    u.PasswordHash = password
}

type Users struct {
    U map[string]*User
}

func Create(name, email string, password []byte) *User {
    return &User{name, email, password, nil}
}

func (u *Users) FindToken(token string) (*User, error) {
    for _, user := range u.U {
        if user.Token != nil && user.Token.Verification == token {
            return user, nil
        }
    }
    return nil, errors.New("Token not found")
}

func (u *Users) Users() []*User {
    users := make([]*User, 0, len(u.U))
    for _, user := range u.U {
        users = append(users, user)
    }
    return users
}

func (u *Users) Add(user *User) error {
    emailHash := Hash(user.Email)
    if _, keyExists := u.U[string(emailHash)]; keyExists {
        return errors.New("User with provided email already exists")
    }
    u.U[string(emailHash)] = user
    return nil
}

func (u *Users) Get(email string) (*User, error) {
    emailHash := Hash(email)
    user, keyExists := u.U[string(emailHash)]
    if !keyExists {
        return nil, errors.New("User with provided email does not exist")
    }
    return user, nil
}

func (u *Users) Delete(user *User) error {
    emailHash := Hash(user.Email)
    if _, keyExists := u.U[string(emailHash)]; !keyExists {
        return errors.New("User with provided email does not exist")
    }
    delete(u.U, string(emailHash))
    return nil
}

type UserDto struct {
    Name string `json:"name"`
    Email string `json:"email"`
    PasswordHash []byte `json:"passwordHash"`
}

func (u *User) ToDto() *UserDto {
    return &UserDto{u.Name, u.Email, u.PasswordHash}
}

func (u *Users) ToDto() []*UserDto {
    users := make([]*UserDto, 0, len(u.U))
    for _, user := range u.U {
        users = append(users, user.ToDto())
    }
    return users
}
