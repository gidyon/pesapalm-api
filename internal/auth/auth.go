package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AccessDetails struct {
	TokenUuid string
	UserId    uint64
	UserName  string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AuthInterface interface {
	CreateAuth(context.Context, uint64, *TokenDetails) error
	FetchAuth(context.Context, string) (string, error)
	DeleteRefresh(context.Context, string) error
	DeleteTokens(context.Context, *AccessDetails) error
}

type RedisAuthService struct {
	client *redis.Client
}

var _ AuthInterface = &RedisAuthService{}

func NewAuthService(client *redis.Client) *RedisAuthService {
	return &RedisAuthService{client: client}
}

// Save token metadata to Redis
func (tk *RedisAuthService) CreateAuth(ctx context.Context, userId uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(ctx, td.TokenUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.client.Set(ctx, td.RefreshUuid, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

// Check the metadata saved
func (tk *RedisAuthService) FetchAuth(ctx context.Context, tokenUuid string) (string, error) {
	userid, err := tk.client.Get(ctx, tokenUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

// Once a user row in the token table
func (tk *RedisAuthService) DeleteTokens(ctx context.Context, authD *AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := tk.client.Del(ctx, authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := tk.client.Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (tk *RedisAuthService) DeleteRefresh(ctx context.Context, refreshUuid string) error {
	//delete refresh token
	deleted, err := tk.client.Del(ctx, refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
