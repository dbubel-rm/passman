package models

const FbJSON = "1"
const LocalID = "2"
const FirebaseResp = "3"

type FirebaseAuthResp struct {
	Kind         string `json:"kind"`
	LocalID      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IDToken      string `json:"idToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	Description  string `json:"description"`
}
