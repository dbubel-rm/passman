package models

type FirebaseSession struct {
	DisplayName  string `json:"displayName"`
	Email        string `json:"email"`
	ExpiresIn    string `json:"expiresIn"`
	IDToken      string `json:"idToken"`
	Kind         string `json:"kind"`
	LocalID      string `json:"localId"`
	RefreshToken string `json:"refreshToken"`
	Registered   bool   `json:"registered"`
}
