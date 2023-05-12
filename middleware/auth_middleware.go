package middleware

import (
	"context"
	"net/http"
	"os"

	"errors"

	"github.com/Prateek61/go_auth/graph/model"
	"github.com/Prateek61/go_auth/postgres"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"currentUser"}

func AuthMiddleware(repo postgres.UsersRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := parseToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}

			user, err := repo.GetUserByID(claims["jti"].(string))
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter: stripBearerPrefx,
}

func stripBearerPrefx(token string) (string, error) {
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:], nil
	}

	return token, errors.New("invalid auth header")
}

var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func parseToken(r *http.Request) (*jwt.Token, error) {
	jwtToken, err := request.ParseFromRequest(r, authExtractor, func(token *jwt.Token) (interface{}, error) {
		t := []byte(os.Getenv("JWT_SECRET"))
		return t, nil
	})

	if err != nil {
		return nil, errors.Join(err, errors.New("token parsing error"))
	}

	return jwtToken, nil
}

func GetCurrentUserFromCTX(ctx context.Context) (*model.User, error) {
	if ctx.Value(userCtxKey) == nil {
		return nil, errors.New("no user in context")
	}

	user, ok := ctx.Value(userCtxKey).(*model.User)
	if !ok || user.ID == "" {
		return nil, errors.New("no user in context")
	}

	return user, nil
}