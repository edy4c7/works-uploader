package services

import (
	"context"
)

//AuthService 認証サービスインターフェース
type AuthService interface {
	//JWTの検証を行う
	VerifyToken(c context.Context, token string) error
}
