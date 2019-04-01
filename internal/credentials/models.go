package credentials

type Add struct {
	ServiceName string `json:"serviceName" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type Credential struct {
	ServiceName string `json:"serviceName" db:"service_name"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
}

type Update struct {
	ServiceName string `json:"serviceName" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type Service struct {
	CredentialID string `json:"credentialId" db:"credential_id"`
	ServiceName  string `json:"serviceName" db:"service_name"`
}
