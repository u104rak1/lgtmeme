package dto

import "github.com/lestrrat-go/jwx/jwk"

type JWKSResp struct {
	Keys []jwk.Key `json:"keys"`
}
