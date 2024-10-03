package models

type Receipt struct {
    ID     string `json:"id"`
    UserID string `json:"user_id"`
    Path   string `json:"path"`
}
