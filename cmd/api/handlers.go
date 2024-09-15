package main

import (
	"authService/internal/models"
	"authService/internal/models/tokens"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var req models.UserReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
	}

	err := app.Users.CreateUser(user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	resp, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (app *application) createTokensHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("user_id")

	ip, err := models.GetIP(r)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := app.Users.GetUser(id)
	if err != nil {
		app.notFound(w, r)
		return
	}

	accessToken, accessClaims, err := app.JWT.CreateToken(id, ip, user.Email, 30*time.Minute)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	refreshTokenHash, refreshToken, refreshClaims, err := app.JWT.CreateRefreshToken(id, ip, user.Email)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	session, err := app.Sessions.CreateSession(&models.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       id,
		RefreshToken: refreshTokenHash,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	resp := tokens.Resp{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	respJson, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respJson)

}

func (app *application) renewAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req tokens.RenewAccessTokenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequest(w, r, err)
		return
	}

	ip, err := models.GetIP(r)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	refreshClaims, err := app.JWT.VerifyToken(req.RefreshToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user, err := app.Users.GetUser(refreshClaims.ID)
	if err != nil {
		app.notFound(w, r)
		return
	}

	log.Println(user)

	if refreshClaims.UserIP != ip {
		err = app.Mailer.Send(user.Email, "warning.gohtml", user.Username)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}

	session, err := app.Sessions.GetSession(refreshClaims.ID)
	if err != nil {
		app.notFound(w, r)
		return
	}

	if session.UserID != refreshClaims.ID {
		app.badRequest(w, r, err)
		return
	}

	accessToken, accessClaims, err := app.JWT.CreateToken(refreshClaims.ID, ip, refreshClaims.Subject, 30*time.Minute)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	resp := tokens.RenewAccessTokenResp{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	respJson, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respJson)
}
