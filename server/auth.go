package server

import (
    "context"
    "fmt"
    "strings"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("segredo-supersecreto")

func JWTUnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        if  info.FullMethod == "/order.OrderService/Login" ||
            info.FullMethod == "/order.OrderService/Register" {
            return handler(ctx, req)
        }


        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, fmt.Errorf("metadata não encontrada")
        }

        authHeader := md["authorization"]
        if len(authHeader) == 0 {
            return nil, fmt.Errorf("token ausente")
        }

        tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
        claims := &jwt.RegisteredClaims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            return nil, fmt.Errorf("token inválido: %v", err)
        }

        return handler(ctx, req)
    }
}

func GenerateToken(username string) (string, error) {
    claims := &jwt.RegisteredClaims{
        Subject: username,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}
