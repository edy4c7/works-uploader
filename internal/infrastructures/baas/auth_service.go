package infrastructures

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
)

//AuthServiceImpl 認証サービスの実装
type AuthServiceImpl struct {
	client *auth.Client
}

//NewAuthServiceImpl 新しい認証サービスのインスタンスを生成する
func NewAuthServiceImpl(client *auth.Client) *AuthServiceImpl {
	return &AuthServiceImpl{
		client: client,
	}
}

//VerifyToken JWTの検証を行う
func (r *AuthServiceImpl) VerifyToken(c context.Context, tokenStr string) error {
	token, err := r.client.VerifyIDToken(c, tokenStr)
	if err != nil {
		return err
	}

	log.Printf("ID token: %v\n", token)
	return nil
}
