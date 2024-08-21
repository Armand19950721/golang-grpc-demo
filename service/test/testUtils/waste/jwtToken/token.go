package token

import (
	"service/model"
	"service/repositories"
	// "service/utils"
	// "context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	// "google.golang.org/grpc/metadata"
)

var (
	jwtKey = []byte("FDr1VjVQiSiybYJrQZNt8Vfd7bFEsKP6vNX1brOSiWl0mAIVCxJiR4/T3zpAlBKc2/9Lw2ac4IwMElGZkssfj3dqwa7CQC7IIB+nVxiM1c9yfowAZw4WQJ86RCUTXaXvRX8JoNYlgXcRrK3BK0E/fKCOY1+izInW3abf0jEeN40HJLkXG6MZnYdhzLnPgLL/TnIFTTAbbItxqWBtkz6FkZTG+dkDSXN7xNUxlg==")
)

type authUser struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

func GenToken(user model.User) (string, time.Time, error) {
	expireTime := time.Now().Add(240 * time.Hour)

	// create token
	tokenUnsign := jwt.NewWithClaims(jwt.SigningMethodHS512, authUser{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Name,
			ExpiresAt: expireTime.Unix(),
		},
		UserId: user.Id.String(),
	})

	token, err := tokenUnsign.SignedString(jwtKey)

	return token, expireTime, err
}

// func CheckIsExpireToken(tokenStr string) bool {
// 	hasKey := utils.HasRedis("expire_token:" + tokenStr)
// 	if hasKey {
// 		utils.PrintObj("found expire token")
// 		return true
// 	}
// 	return false
// }

func ValidToken(tokenString string) (*model.User, error) {

	// isExpire := CheckIsExpireToken(tokenString)
	// if isExpire {
	// 	return nil, errors.New("expire token")
	// }

	var claims authUser
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// utils.PrintObj(claims, "claims")

	userQuery := model.User{Id: uuid.MustParse(claims.UserId)}
	resultQuery, userData := repositories.QueryUser(userQuery)

	if resultQuery.Error != nil {
		return nil, errors.New("invalid user")
	}

	return &userData, nil
}

// func ValidTokenByCtx(ctx context.Context) (token string, valid bool) {
// 	// parse metadata
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		utils.PrintObj("md", "get basic info md err")
// 		return "", false
// 	}

// 	// check token
// 	tokenStr := utils.GetMetaDataField(md, "token")
// 	_, err := ValidToken(tokenStr)

// 	if err != nil {
// 		return "", false
// 	}

// 	return tokenStr, true //if no err then success
// }
