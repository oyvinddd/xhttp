package jsrouter

import (
	"context"
	"net/http"
	"github.com/oyvinddd/xhttp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func Authorize(next httprouter.Handle, requiredRole xhttp.Role) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenStr, err := xhttp.ParseAccessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessClaims := new(xhttp.AccessTokenClaims)

		token, err := jwt.ParseWithClaims(tokenStr, accessClaims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, xhttp.UnexpectedSigningError
			}
			return xhttp.GetHMACSecret(), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, xhttp.InvalidAccessTokenError.Error(), http.StatusUnauthorized)
			return
		}

		if !accessClaims.Role.HasPermission(requiredRole) {
			http.Error(w, xhttp.InsufficientRoleError.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), xhttp.AccessClaimsCtxKey, accessClaims)
		next(w, r.WithContext(ctx), ps)
	}
}

