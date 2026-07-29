package httprouter

import (
	"net/http"
	"github.com/oyvinddd/xtoken"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(manager xtoken.Manager, next httprouter.Handle, requiredRole Role) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenStr, err := xtoken.ParseAccessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessClaims := new(AccessTokenClaims)

		token, err := jwt.ParseWithClaims(tokenStr, accessClaims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, xtoken.UnexpectedSigningError
			}
			return manager.secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, xtoken.InvalidAccessTokenError.Error(), http.StatusUnauthorized)
			return
		}

		if !accessClaims.Role.HasPermission(requiredRole) {
			http.Error(w, xtoken.InsufficientRoleError.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), xtoken.AccessClaimsCtxKey, accessClaims)
		next(w, r.WithContext(ctx), ps)
	}
}

