package credentials

type Credential struct {
	CredentialID string `json:"credentialId" db:"credential_id"`
	LocalID      string `json:"localId" db:"local_id"`
	ServiceName  string `json:"serviceName" db:"service_name"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	UpdatedAt    string `json:"updatedAt" db:"updated_at"`
	CreatedAt    string `json:"createdAt" db:"created_at"`
}
