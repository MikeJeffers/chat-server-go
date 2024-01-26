package main

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

type CustomClaims struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"name"`
}

func redisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
		DB:       0,
		Protocol: 3,
	})
}

func verifyToken(rawJwt string, redisClient *redis.Client) *User {
	user := validateJwt(rawJwt)
	if user == nil {
		fmt.Println("validate jwt failed")
		return nil
	}
	value, err := redisClient.Get(context.TODO(), fmt.Sprintf("jwt:%v", user.Id)).Bytes()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Verify the provided token matches what was stored by the auth api
	redisJwt := string(value)
	if rawJwt != redisJwt {
		fmt.Printf("tokens dont match! %s != %s  \n", rawJwt, redisJwt)
		return nil
	}
	return user
}

func validateJwt(tokenString string) *User {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return ([]byte(SECRET_JWT)), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	mapped := token.Claims.(jwt.MapClaims)

	if mapped["id"] == nil || mapped["username"] == nil {
		return nil
	}
	return &User{int64(mapped["id"].(float64)), mapped["username"].(string)}
}
