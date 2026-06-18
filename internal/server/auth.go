package server

import (
"encoding/json"
"errors"
"fmt"
"io"
"log"
"net/http"
"net/url"
"os"
"strings"
"time"

"github.com/SorynMochi/Namagotchi/internal/database"
)

const authSessionCookieName = "namigotchi_session"

type AuthRegisterRequest struct {
DisplayName string `json:"displayName"`
Email       string `json:"email"`
Password    string `json:"password"`
}

type AuthLoginRequest struct {
Email    string `json:"email"`
Password string `json:"password"`
}

type AuthResponse struct {
OK       bool                  `json:"ok"`
LoggedIn bool                  `json:"loggedIn"`
Account  *database.AuthAccount `json:"account,omitempty"`
Message  string                `json:"message"`
}

type googleTokenResponse struct {
AccessToken string `json:"access_token"`
IDToken     string `json:"id_token"`
TokenType   string `json:"token_type"`
ExpiresIn   int    `json:"expires_in"`
Error       string `json:"error"`
ErrorDesc   string `json:"error_description"`
}

type googleTokenInfoResponse struct {
Aud           string `json:"aud"`
Sub           string `json:"sub"`
Email         string `json:"email"`
EmailVerified string `json:"email_verified"`
Name          string `json:"name"`
Picture       string `json:"picture"`
}

func (s *Server) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
writeError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

var request AuthRegisterRequest
if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
writeError(w, http.StatusBadRequest, "invalid registration request")
return
}

request.DisplayName = strings.TrimSpace(request.DisplayName)
request.Email = strings.TrimSpace(request.Email)

if len(request.DisplayName) < 2 || len(request.DisplayName) > 32 {
writeError(w, http.StatusBadRequest, "display name must be 2 to 32 characters")
return
}

if !looksLikeEmail(request.Email) {
writeError(w, http.StatusBadRequest, "email is invalid")
return
}

if len(request.Password) < 8 {
writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
return
}

account, sessionToken, err := s.Store.RegisterGameAccount(
r.Context(),
request.DisplayName,
request.Email,
request.Password,
)
if err != nil {
switch {
case errors.Is(err, database.ErrAuthDisplayNameTaken):
writeError(w, http.StatusConflict, "display name is already taken")
case errors.Is(err, database.ErrAuthEmailTaken):
writeError(w, http.StatusConflict, "email is already registered")
default:
log.Printf("register auth account failed: %v", err)
writeError(w, http.StatusInternalServerError, "registration failed")
}
return
}

setAuthSessionCookie(w, r, sessionToken)
writeJSON(w, http.StatusOK, AuthResponse{
OK:       true,
LoggedIn: true,
Account:  &account,
Message:  "Account created.",
})
}

func (s *Server) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
writeError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

var request AuthLoginRequest
if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
writeError(w, http.StatusBadRequest, "invalid login request")
return
}

account, sessionToken, err := s.Store.LoginGameAccount(
r.Context(),
request.Email,
request.Password,
)
if err != nil {
if errors.Is(err, database.ErrAuthInvalidCredentials) {
writeError(w, http.StatusUnauthorized, "invalid email or password")
return
}

log.Printf("login failed: %v", err)
writeError(w, http.StatusInternalServerError, "login failed")
return
}

setAuthSessionCookie(w, r, sessionToken)
writeJSON(w, http.StatusOK, AuthResponse{
OK:       true,
LoggedIn: true,
Account:  &account,
Message:  "Logged in.",
})
}

func (s *Server) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost && r.Method != http.MethodGet {
writeError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

if token := authSessionTokenFromRequest(r); token != "" {
if err := s.Store.DeleteAuthSession(r.Context(), token); err != nil {
log.Printf("delete auth session failed: %v", err)
}
}

clearAuthSessionCookie(w, r)
writeJSON(w, http.StatusOK, AuthResponse{
OK:       true,
LoggedIn: false,
Message:  "Logged out.",
})
}

func (s *Server) HandleAuthMe(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodGet {
writeError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

account, ok := s.AuthAccountFromRequest(r)
if !ok {
writeJSON(w, http.StatusOK, AuthResponse{
OK:       true,
LoggedIn: false,
Message:  "Not logged in.",
})
return
}

writeJSON(w, http.StatusOK, AuthResponse{
OK:       true,
LoggedIn: true,
Account:  &account,
Message:  "Logged in.",
})
}

func (s *Server) HandleAuthGoogleStart(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodGet {
writeError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

clientID, _, redirectURI, ok := googleOAuthConfig(r)
if !ok {
writeError(w, http.StatusServiceUnavailable, "google login is not configured")
return
}

redirectPath := r.URL.Query().Get("redirect")
state, err := s.Store.CreateOAuthState(r.Context(), "google", redirectPath)
if err != nil {
log.Printf("create google oauth state failed: %v", err)
writeError(w, http.StatusInternalServerError, "google login failed")
return
}

params := url.Values{}
params.Set("client_id", clientID)
params.Set("redirect_uri", redirectURI)
params.Set("response_type", "code")
params.Set("scope", "openid email profile")
params.Set("state", state)
params.Set("prompt", "select_account")

http.Redirect(w, r, "https://accounts.google.com/o/oauth2/v2/auth?"+params.Encode(), http.StatusFound)
}

func (s *Server) HandleAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodGet {
writeError(w, http.StatusMethodNotAllowed, "method not allowed")
return
}

errorParam := strings.TrimSpace(r.URL.Query().Get("error"))
if errorParam != "" {
writeError(w, http.StatusBadRequest, "google login was cancelled or denied")
return
}

code := strings.TrimSpace(r.URL.Query().Get("code"))
state := strings.TrimSpace(r.URL.Query().Get("state"))
if code == "" || state == "" {
writeError(w, http.StatusBadRequest, "missing google login code")
return
}

clientID, clientSecret, redirectURI, ok := googleOAuthConfig(r)
if !ok {
writeError(w, http.StatusServiceUnavailable, "google login is not configured")
return
}

redirectPath, err := s.Store.ConsumeOAuthState(r.Context(), "google", state)
if err != nil {
if errors.Is(err, database.ErrAuthInvalidState) {
writeError(w, http.StatusBadRequest, "google login state expired")
return
}

log.Printf("consume google oauth state failed: %v", err)
writeError(w, http.StatusInternalServerError, "google login failed")
return
}

tokenResponse, err := exchangeGoogleCode(code, clientID, clientSecret, redirectURI)
if err != nil {
log.Printf("exchange google code failed: %v", err)
writeError(w, http.StatusBadGateway, "google login failed")
return
}

tokenInfo, err := fetchGoogleTokenInfo(tokenResponse.IDToken)
if err != nil {
log.Printf("fetch google token info failed: %v", err)
writeError(w, http.StatusBadGateway, "google login failed")
return
}

if tokenInfo.Aud != clientID || tokenInfo.Sub == "" {
writeError(w, http.StatusUnauthorized, "google login verification failed")
return
}

account, sessionToken, err := s.Store.FindOrCreateExternalAuthAccount(
r.Context(),
"google",
tokenInfo.Sub,
tokenInfo.Name,
tokenInfo.Email,
tokenInfo.Picture,
)
if err != nil {
log.Printf("find/create google auth account failed: %v", err)
writeError(w, http.StatusInternalServerError, "google login failed")
return
}

setAuthSessionCookie(w, r, sessionToken)
log.Printf("google login succeeded for account %d (%s)", account.ID, account.DisplayName)

http.Redirect(w, r, redirectPath, http.StatusFound)
}

func (s *Server) AuthAccountFromRequest(r *http.Request) (database.AuthAccount, bool) {
token := authSessionTokenFromRequest(r)
if token == "" {
return database.AuthAccount{}, false
}

account, err := s.Store.AccountByAuthSession(r.Context(), token)
if err != nil {
return database.AuthAccount{}, false
}

return account, true
}

func exchangeGoogleCode(code, clientID, clientSecret, redirectURI string) (googleTokenResponse, error) {
var tokenResponse googleTokenResponse

form := url.Values{}
form.Set("code", code)
form.Set("client_id", clientID)
form.Set("client_secret", clientSecret)
form.Set("redirect_uri", redirectURI)
form.Set("grant_type", "authorization_code")

request, err := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
if err != nil {
return tokenResponse, err
}

request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

response, err := http.DefaultClient.Do(request)
if err != nil {
return tokenResponse, err
}
defer response.Body.Close()

body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
if err != nil {
return tokenResponse, err
}

if err := json.Unmarshal(body, &tokenResponse); err != nil {
return tokenResponse, err
}

if response.StatusCode < 200 || response.StatusCode >= 300 {
if tokenResponse.ErrorDesc != "" {
return tokenResponse, errors.New(tokenResponse.ErrorDesc)
}
if tokenResponse.Error != "" {
return tokenResponse, errors.New(tokenResponse.Error)
}

return tokenResponse, fmt.Errorf("google token exchange failed: %s", response.Status)
}

if tokenResponse.IDToken == "" {
return tokenResponse, errors.New("google token response did not include id_token")
}

return tokenResponse, nil
}

func fetchGoogleTokenInfo(idToken string) (googleTokenInfoResponse, error) {
var tokenInfo googleTokenInfoResponse

if strings.TrimSpace(idToken) == "" {
return tokenInfo, errors.New("missing id token")
}

endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
response, err := http.Get(endpoint)
if err != nil {
return tokenInfo, err
}
defer response.Body.Close()

body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
if err != nil {
return tokenInfo, err
}

if err := json.Unmarshal(body, &tokenInfo); err != nil {
return tokenInfo, err
}

if response.StatusCode < 200 || response.StatusCode >= 300 {
return tokenInfo, fmt.Errorf("google tokeninfo failed: %s", response.Status)
}

return tokenInfo, nil
}

func googleOAuthConfig(r *http.Request) (clientID, clientSecret, redirectURI string, ok bool) {
clientID = strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
clientSecret = strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
redirectURI = strings.TrimSpace(os.Getenv("GOOGLE_REDIRECT_URL"))

if redirectURI == "" {
baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("APP_PUBLIC_URL")), "/")
if baseURL != "" {
redirectURI = baseURL + "/api/auth/google/callback"
} else if r != nil {
scheme := "http"
if r.TLS != nil {
scheme = "https"
}

redirectURI = fmt.Sprintf("%s://%s/api/auth/google/callback", scheme, r.Host)
}
}

return clientID, clientSecret, redirectURI, clientID != "" && clientSecret != "" && redirectURI != ""
}

func authSessionTokenFromRequest(r *http.Request) string {
cookie, err := r.Cookie(authSessionCookieName)
if err != nil {
return ""
}

return strings.TrimSpace(cookie.Value)
}

func setAuthSessionCookie(w http.ResponseWriter, r *http.Request, token string) {
http.SetCookie(w, &http.Cookie{
Name:     authSessionCookieName,
Value:    token,
Path:     "/",
MaxAge:   int((30 * 24 * time.Hour).Seconds()),
Expires:  time.Now().Add(30 * 24 * time.Hour),
HttpOnly: true,
SameSite: http.SameSiteLaxMode,
Secure:   authCookieShouldBeSecure(r),
})
}

func clearAuthSessionCookie(w http.ResponseWriter, r *http.Request) {
http.SetCookie(w, &http.Cookie{
Name:     authSessionCookieName,
Value:    "",
Path:     "/",
MaxAge:   -1,
Expires:  time.Unix(0, 0),
HttpOnly: true,
SameSite: http.SameSiteLaxMode,
Secure:   authCookieShouldBeSecure(r),
})
}

func authCookieShouldBeSecure(r *http.Request) bool {
if strings.TrimSpace(os.Getenv("AUTH_SECURE_COOKIE")) == "1" {
return true
}

return r != nil && r.TLS != nil
}

func looksLikeEmail(email string) bool {
email = strings.TrimSpace(email)
return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) <= 254
}
