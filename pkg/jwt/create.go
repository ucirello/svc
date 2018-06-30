// Package jwt manipulates JWT tokens.
package jwt // import "cirello.io/svc/pkg/jwt"

import (
	"crypto/x509"

	"cirello.io/errors"
	jwt "github.com/dgrijalva/jwt-go"
)

// ServiceClaims define the set of claims used by cirello.io services.
type ServiceClaims struct {
	jwt.StandardClaims

	// Email is the actor who is logging in.
	Email string
	// Target defines to which service this token was created for.
	Target string
	// Trust defines the trust level so to give the application some context
	// on how it should handle low-trust logins.
	Trust string
}

// CreateFromCert a JWT whose content indicate a high-trust login.
func CreateFromCert(svcName string, caPEM []byte, cert *x509.Certificate) (string, error) {
	if len(cert.EmailAddresses) == 0 {
		return "", errors.E("certificate missing email")
	} else if len(cert.EmailAddresses) > 1 {
		return "", errors.E("multiple emails in the same certificate - cannot choose one")
	}

	email := cert.EmailAddresses[0]
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&ServiceClaims{
			Email:  email,
			Target: svcName,
			Trust:  "high",
		},
	)

	tokenString, err := token.SignedString(caPEM)
	return tokenString, errors.E(err, "cannot sign JWT")
}

// CreateFromEmail a JWT whose content indicate a low-trust login.
func CreateFromEmail(svcName string, caPEM []byte, email string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&ServiceClaims{
			Email:  email,
			Target: svcName,
			Trust:  "low",
		},
	)

	tokenString, err := token.SignedString(caPEM)
	return tokenString, errors.E(err, "cannot sign JWT")
}