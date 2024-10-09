package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenManager struct {
	*TokenOptions
}

type TokenOptions struct {
	AccessSecret      string
	RefreshSecret     string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

func NewTokenService(opt *TokenOptions) *TokenManager {
	return &TokenManager{TokenOptions: opt}
}

type TokenInterface interface {
	CreateToken(userId uint64, userName string) (*TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
	TokenValid(*http.Request) error
	VerifyToken(*http.Request) (*jwt.Token, error)
}

// Token implements the TokenInterface
var _ TokenInterface = &TokenManager{}

func (t *TokenManager) CreateToken(userId uint64, userName string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(t.AccessExpiration).Unix() //expires after 30 min
	td.TokenUuid = uuid.NewString()

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["user_id"] = userId
	atClaims["user_name"] = userName
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(t.AccessSecret))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	td.RtExpires = time.Now().Add(t.RefreshExpiration).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + fmt.Sprint(userId)

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["user_name"] = userName
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(t.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (t *TokenManager) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := t.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	acc, err := Extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (t *TokenManager) TokenValid(r *http.Request) error {
	token, err := t.VerifyToken(r)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

func (t *TokenManager) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if r.URL.Path == "/api/refresh" {
			return []byte(t.RefreshSecret), nil
		}
		return []byte(t.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// get the token from the request body
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func Extract(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok1 := claims["access_uuid"].(string)
		refreshUuid, ok2 := claims["refresh_uuid"].(string)
		userId, userOk := claims["user_id"].(float64)
		userName, userNameOk := claims["user_name"].(string)
		if (!ok1 && !ok2) || !userOk || !userNameOk {
			return nil, errors.New("unauthorized")
		} else {
			return &AccessDetails{
				TokenUuid: firstVal(refreshUuid, accessUuid),
				UserId:    uint64(userId),
				UserName:  userName,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

func firstVal(vals ...string) string {
	for _, str := range vals {
		if str != "" {
			return str
		}
	}
	return ""
}
