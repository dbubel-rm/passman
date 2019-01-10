package models

type Credentials struct {
	// CredentialID int       `json:"credentialId" db:"credential_id"`
	// UserID       int       `json:"userId" db:"user_id"`
	ServiceName string `json:"serviceName" db:"service_name"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	// UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	// CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	// DeletedAt    time.Time `json:"deletedAt" db:"deleted_at"`
}
