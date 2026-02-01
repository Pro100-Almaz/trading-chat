package controller

import (
	"net/http"

	"github.com/Pro100-Almaz/trading-chat/bootstrap"
	"github.com/Pro100-Almaz/trading-chat/domain"
	"github.com/Pro100-Almaz/trading-chat/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type GoogleController struct {
	GoogleUseCase domain.GoogleUseCase
	Env           *bootstrap.Env
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL: "http://localhost:8080/api/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

// HandleGoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects user to Google OAuth consent screen
// @Tags Authentication
// @Success 307 "Redirect to Google"
// @Router /google/login [get]
func (gc *GoogleController) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := gc.GoogleUseCase.GenerateStateOauthCookie(w)
	googleOauthConfig.ClientSecret = gc.Env.GoogleClientSecret
	googleOauthConfig.ClientID = gc.Env.GoogleClientID
	u := googleOauthConfig.AuthCodeURL(oauthState)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handles the callback from Google OAuth and creates/logs in user
// @Tags Authentication
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter for CSRF protection"
// @Success 307 "Redirect to profile page with auth cookies"
// @Failure 307 "Redirect to home on error"
// @Router /google/callback [get]
func (gc *GoogleController) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	googleOauthConfig.ClientSecret = gc.Env.GoogleClientSecret
	googleOauthConfig.ClientID = gc.Env.GoogleClientID
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Error("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := gc.GoogleUseCase.GetUserDataFromGoogle(googleOauthConfig, r.FormValue("code"), oauthGoogleUrlAPI)
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	accessToken, refreshToken, err := gc.GoogleUseCase.GoogleLogin(ctx, data, gc.Env)
	if err != nil {
		log.Error(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// write access token and refresh token to cookie
	utils.SetCookie(w, "access_token", accessToken)
	utils.SetCookie(w, "refresh_token", refreshToken)

	// redirect to home page
	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}
