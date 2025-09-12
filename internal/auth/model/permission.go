package model

import "time"

type Permission struct {
    UserID    int       `json:"user_id"`
    Module    string    `json:"module"`
    Action    string    `json:"action"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
