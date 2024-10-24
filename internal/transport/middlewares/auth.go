package middlewares

import (
	"context"
	"net/http"
	"seeker/internal/domain/dto"
	errs "seeker/internal/domain/errors"
	"seeker/pkg/handler/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func WithAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

		tokens, err := getAuthToken(r)

		if err != nil {
			response.Error(w, err, http.StatusForbidden)
			return
		}

		session := &dto.JWTSession{}

		err = parseJWTSession(tokens.AccessToken, session)

		if err != nil {
			response.Error(w, err, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), dto.CtxSessionKey, &session.User)
		r = r.WithContext(ctx)

		h(w, r, params)
	}
}

func getAuthToken(r *http.Request) (dto.JWTTokenResponse, error) {

	var tokens dto.JWTTokenResponse

	accessToken, err := r.Cookie(dto.AccessTokenCookieKey)
	refreshToken, err := r.Cookie(dto.RefreshTokenCookieKey)

	if err != nil {
		return tokens, errs.ErrUnauthorized
	}

	tokens.AccessToken = accessToken.Value
	tokens.RefreshToken = refreshToken.Name

	return tokens, nil
}

func parseJWTSession(token string, session *dto.JWTSession) error {
	_, err := jwt.ParseWithClaims(token, session, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return errs.ErrFailedToParseJWTClaims
	}

	return nil
}
