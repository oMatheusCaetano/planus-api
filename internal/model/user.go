package model

import "time"

type User struct {
	ID        ID        `bson:"_id"        json:"id"`
	Name      string    `bson:"name"       json:"name"`
	Username  string    `bson:"username"   json:"username"`
	Password  string    `bson:"password"   json:"-"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func NewUser(name, username, hashedPassword string) *User {
    now := time.Now()

    return &User{
        ID:        NewID(),
        Name:      name,
        Username:  username,
        Password:  hashedPassword,
        CreatedAt: now,
        UpdatedAt: now,
    }
}
