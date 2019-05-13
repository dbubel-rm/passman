package models

type Credential struct {
	CredentialID string `json:"credential_id" db:"credential_id"`
	LocalID      string `json:"local_id" db:"local_id"`
	ServiceName  string `json:"service_name" db:"service_name"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	UpdatedAt    string `json:"updated_at" db:"updated_at"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}
