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


// contextKey is a type for custom context keys 
type contextKey struct {
	name string
}
var userCtxKey = &contextKey{"currentUser"} // context key for the current user

// AuthMiddleware is a middleware for authentication
func AuthMiddleware(repo postgres.UsersRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
			// Parse JWT token
			token, err := parseToken(req)
			if err != nil {
				next.ServeHTTP(responseWriter, req)
				return
			}

			// Get claims from token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				next.ServeHTTP(responseWriter, req)
				return
			}

			// Get user from DB
			user, err := repo.GetUserByID(claims["jti"].(string))
			if err != nil {
				next.ServeHTTP(responseWriter, req)
				return
			}

			// Add user to context
			ctx := context.WithValue(req.Context(), userCtxKey, user)

			// Call next with our new context
			next.ServeHTTP(responseWriter, req.WithContext(ctx))
		})
	}
}

// Extracts the token from the Authorization header
var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter: stripBearerPrefx,
}

// Strips the Bearer prefix from the token
func stripBearerPrefx(token string) (string, error) {
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:], nil
	}

	return token, errors.New("invalid auth header")
}

// Extracts the token from the Authorization header or from the access_token argument
var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

// Parses the token from the request
func parseToken(req *http.Request) (*jwt.Token, error) {
	// Parse token from request
	jwtToken, err := request.ParseFromRequest(req, authExtractor, func(token *jwt.Token) (interface{}, error) {
		t := []byte(os.Getenv("JWT_SECRET"))
		return t, nil
	})
	if err != nil {
		return nil, errors.Join(err, errors.New("token parsing error"))
	}

	return jwtToken, nil
}

// GetCurrentUserFromCTX gets the current user from the context
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