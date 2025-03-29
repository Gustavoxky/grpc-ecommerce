package server

import (
    "context"
    "log"
    "strings"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"

    "github.com/golang-jwt/jwt/v5"
)

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        start := time.Now()

        // Tenta extrair o usuário do token JWT
        user := "anônimo"
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            auth := md.Get("authorization")
            if len(auth) > 0 && strings.HasPrefix(auth[0], "Bearer ") {
                tokenStr := strings.TrimPrefix(auth[0], "Bearer ")
                claims := &jwt.RegisteredClaims{}
                token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
                    return jwtKey, nil
                })
                if err == nil && token.Valid {
                    user = claims.Subject
                }
            }
        }

        // Executa o método real
        res, err := handler(ctx, req)

        duration := time.Since(start)
        status := "OK"
        if err != nil {
            status = "ERRO: " + err.Error()
        }

        log.Printf("[gRPC] método: %s | usuário: %s | duração: %s | status: %s",
            info.FullMethod, user, duration.String(), status,
        )

        return res, err
    }
}
