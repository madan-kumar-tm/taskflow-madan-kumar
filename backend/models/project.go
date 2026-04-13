package models

type Project struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	OwnerID     string  `json:"owner_id"`
	CreatedAt   string  `json:"created_at"`
}
