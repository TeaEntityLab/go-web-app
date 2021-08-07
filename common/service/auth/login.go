package serviceAuth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	bcrypt "golang.org/x/crypto/bcrypt"

	mod "go-web-app/common/model"
	errtrace "go-web-app/errtrace"
)

// CheckAuthInfoValidation ...
func CheckAuthInfoValidation(auth *mod.AuthLogin) bool {
	return auth.UserName != "" && auth.Password != ""
}

// CheckAuthTokenValidation ...
func CheckAuthTokenValidation(authTokenJWTString string) (*mod.AuthToken, *errtrace.Error) {
	if len(authTokenJWTString) == 0 {
		return nil, &errtrace.Error{
			ErrCode: ErrorAuthTokenInvalidEmpty,
		}
	}
	authToken, decodeErr := decodeJWTTokenForUser(authTokenJWTString)
	if authToken == nil || decodeErr != nil {
		return nil, &errtrace.Error{
			ErrRef:  decodeErr,
			ErrCode: ErrorAuthTokenInvalidDecodeError,
		}
	}
	if time.Now().UTC().UnixNano()-authToken.Ttl > IntervalTtlToken {
		return nil, &errtrace.Error{
			ErrCode: ErrorAuthTokenInvalidTtlTimeout,
		}
	}

	return authToken, nil
}

// CheckLoginUserNamePassword ...
func CheckLoginUserNamePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateLoginHashedPassword ...
func GenerateLoginHashedPassword(password string) (string, error) {
	hashedPasswordByteArray, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPasswordByteArray), err
}

// GenerateJWTTokenForUser ...
func GenerateJWTTokenForUser(user *mod.User) (string, error) {
	return GenerateJWTTokenByUserNameAndUserID(user.UserName, user.ModelID())
}

// GenerateJWTTokenByUserNameAndUserID ...
func GenerateJWTTokenByUserNameAndUserID(username string, userID string) (string, error) {
	auth := mod.AuthToken{
		StandardClaims: &jwt.StandardClaims{},
		UserName:       username,
		UserID:         userID,
		Ttl:            time.Now().UTC().UnixNano() + IntervalTtlToken,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, auth)

	return token.SignedString(EnvJWTTokenPrivateKey)
}

// decodeJWTTokenForUser ...
func decodeJWTTokenForUser(tokenString string) (*mod.AuthToken, error) {
	authToken := &mod.AuthToken{
		StandardClaims: &jwt.StandardClaims{},
	}
	_, err := jwt.ParseWithClaims(tokenString, authToken, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return EnvJWTTokenPublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return authToken, nil
}
