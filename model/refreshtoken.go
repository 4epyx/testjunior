package model

// RefreshToken represents a refresh token data, stored in database
type RefreshToken struct {
	//Id is a token UUID
	Id string `json:"id" bson:"_id"`
	//UserGuid is a token GUID
	UserGuid string `json:"user_guid" bson:"user_guid"`
	//Expiration is a time of expiration in unix format
	Expiration int64 `json:"exp" bson:"exp"`
	//Token is a hash of the refresh token
	Token string `json:"token" bson:"token"`
}
