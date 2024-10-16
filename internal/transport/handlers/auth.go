package handlers

import (
	"net/http"
	"seeker/internal/domain/dto"
	"seeker/internal/domain/usecases"
	"seeker/internal/transport/middlewares"
	"seeker/internal/types"
	"seeker/pkg/handler"
	"seeker/pkg/handler/request"
	"seeker/pkg/handler/response"

	"github.com/julienschmidt/httprouter"
)

type authHandler struct {
	usecase usecases.AuthUsecase
}

func NewAuthHandler(usecase usecases.AuthUsecase) handler.Handler {
	return &authHandler{
		usecase,
	}
}

const (
	register    = "/auth/register"
	login       = "/auth/login"
	verifyEmail = "/auth/verify-email"
)

func (h *authHandler) Register(router *httprouter.Router) {
	router.POST(register, h.handleRegister)
	router.POST(login, h.handleLogin)
	router.GET(verifyEmail, middlewares.WithAuth(h.handleVerifyEmail))
}

func (h *authHandler) handleRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body dto.RegisterUserInput

	if err := request.ReadBody(r, &body); err != nil {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	input := dto.RegisterUserInput{
		Email:    body.Email,
		Password: body.Password,
	}

	tokens, session, err := h.usecase.Register(input)

	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	response.PrivateCookie(w, types.AccessTokenCookieKey, tokens.AccessToken)
	response.PrivateCookie(w, types.RefreshTokenCookieKey, tokens.RefreshToken)
	response.JSON(w, session, http.StatusCreated)
}

func (h *authHandler) handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body dto.RegisterUserInput

	if err := request.ReadBody(r, &body); err != nil {
		response.Error(w, nil, http.StatusBadRequest)
		return
	}

	input := dto.LoginUserInput{
		Email:    body.Email,
		Password: body.Password,
	}

	tokens, session, err := h.usecase.Login(input)

	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	response.PrivateCookie(w, types.AccessTokenCookieKey, tokens.AccessToken)
	response.PrivateCookie(w, types.RefreshTokenCookieKey, tokens.RefreshToken)
	response.JSON(w, session, http.StatusOK)
}

func (h *authHandler) handleVerifyEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := request.GetSession(r)

	if err != nil {
		response.Error(w, err, http.StatusForbidden)
		return
	}

	tokens, newSession, err := h.usecase.VerifyEmail(session.Email)

	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	response.PrivateCookie(w, types.AccessTokenCookieKey, tokens.AccessToken)
	response.PrivateCookie(w, types.RefreshTokenCookieKey, tokens.RefreshToken)
	response.JSON(w, newSession, http.StatusOK)
}
