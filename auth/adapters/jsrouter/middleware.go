package jsrouter

import (
	"context"
	"net/http"
	auth "github.com/oyvinddd/xhttp/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func Authorize(next httprouter.Handle, requiredRole auth.Role) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenStr, err := auth.ParseAccessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessClaims := new(auth.AccessTokenClaims)

		token, err := jwt.ParseWithClaims(tokenStr, accessClaims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrUnexpectedSigning
			}
			return auth.GetHMACSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, auth.ErrInvalidAccessToken.Error(), http.StatusUnauthorized)
			return
		}

		if !accessClaims.Role.HasPermission(requiredRole) {
			http.Error(w, auth.ErrInsufficientRole.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.AccessClaimsCtxKey, accessClaims)
		next(w, r.WithContext(ctx), ps)
	}
}

