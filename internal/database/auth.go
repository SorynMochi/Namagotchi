package database

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAuthEmailTaken          = errors.New("auth email already registered")
	ErrAuthDisplayNameTaken    = errors.New("auth display name already registered")
	ErrAuthDisplayNameReserved = errors.New("auth display name is reserved")
	ErrAuthInvalidCredentials  = errors.New("invalid auth credentials")
	ErrAuthInvalidSession      = errors.New("invalid auth session")
	ErrAuthInvalidState        = errors.New("invalid oauth state")
)

const (
	authSessionDuration = 30 * 24 * time.Hour
	authStateDuration   = 10 * time.Minute
)

type AuthAccount struct {
	ID          int64      `json:"id"`
	DisplayName string     `json:"displayName"`
	Email       string     `json:"email"`
	AvatarURL   string     `json:"avatarUrl"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

func (s *Store) RegisterGameAccount(ctx context.Context, displayName, email, password string) (AuthAccount, string, error) {
	var account AuthAccount

	displayName = cleanAuthDisplayName(displayName)
	email = normalizeAuthEmail(email)

	if err := ValidateAuthDisplayName(displayName); err != nil {
		return account, "", err
	}

	if email == "" || password == "" {
		return account, "", ErrAuthInvalidCredentials
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return account, "", err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return account, "", err
	}
	defer tx.Rollback(ctx)

	var existingID int64
	err = tx.QueryRow(ctx, `
select id
from auth_accounts
where lower(display_name) = lower($1)
`, displayName).Scan(&existingID)
	if err == nil {
		return account, "", ErrAuthDisplayNameTaken
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return account, "", err
	}

	err = tx.QueryRow(ctx, `
select id
from auth_accounts
where lower(email) = lower($1)
and email <> ''
`, email).Scan(&existingID)
	if err == nil {
		return account, "", ErrAuthEmailTaken
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return account, "", err
	}

	err = tx.QueryRow(ctx, `
insert into auth_accounts (display_name, email)
values ($1, $2)
returning id, display_name, email, avatar_url, created_at, updated_at, last_login_at
`, displayName, email).Scan(
		&account.ID,
		&account.DisplayName,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
	)
	if err != nil {
		return account, "", err
	}

	_, err = tx.Exec(ctx, `
insert into auth_credentials (account_id, email_normalized, password_hash)
values ($1, $2, $3)
`, account.ID, email, string(passwordHash))
	if err != nil {
		return account, "", err
	}

	_, err = tx.Exec(ctx, `
insert into auth_identities (account_id, provider, provider_user_id, email, display_name)
values ($1, 'game', $2, $3, $4)
`, account.ID, email, email, account.DisplayName)
	if err != nil {
		return account, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return account, "", err
	}

	token, err := s.CreateAuthSession(ctx, account.ID)
	if err != nil {
		return account, "", err
	}

	return account, token, nil
}

func (s *Store) LoginGameAccount(ctx context.Context, email, password string) (AuthAccount, string, error) {
	var account AuthAccount
	var passwordHash string

	email = normalizeAuthEmail(email)

	if email == "" || password == "" {
		return account, "", ErrAuthInvalidCredentials
	}

	err := s.Pool.QueryRow(ctx, `
select a.id, a.display_name, a.email, a.avatar_url, a.created_at, a.updated_at, a.last_login_at, c.password_hash
from auth_credentials c
join auth_accounts a on a.id = c.account_id
where c.email_normalized = $1
`, email).Scan(
		&account.ID,
		&account.DisplayName,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
		&passwordHash,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return account, "", ErrAuthInvalidCredentials
	}
	if err != nil {
		return account, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return account, "", ErrAuthInvalidCredentials
	}

	token, err := s.CreateAuthSession(ctx, account.ID)
	if err != nil {
		return account, "", err
	}

	if err := s.touchAuthLogin(ctx, account.ID, "game", email); err != nil {
		return account, "", err
	}

	return account, token, nil
}

func (s *Store) FindOrCreateExternalAuthAccount(ctx context.Context, provider, providerUserID, displayName, email, avatarURL string) (AuthAccount, string, error) {
	var account AuthAccount

	provider = strings.TrimSpace(strings.ToLower(provider))
	providerUserID = strings.TrimSpace(providerUserID)
	_ = cleanAuthDisplayName(displayName)
	displayName = "NamiFan"
	email = normalizeAuthEmail(email)
	avatarURL = strings.TrimSpace(avatarURL)

	if provider == "" || providerUserID == "" {
		return account, "", ErrAuthInvalidCredentials
	}

	err := s.Pool.QueryRow(ctx, `
select a.id, a.display_name, a.email, a.avatar_url, a.created_at, a.updated_at, a.last_login_at
from auth_identities i
join auth_accounts a on a.id = i.account_id
where i.provider = $1 and i.provider_user_id = $2
`, provider, providerUserID).Scan(
		&account.ID,
		&account.DisplayName,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
	)
	if err == nil {
		token, sessionErr := s.CreateAuthSession(ctx, account.ID)
		if sessionErr != nil {
			return account, "", sessionErr
		}

		if loginErr := s.touchAuthLogin(ctx, account.ID, provider, providerUserID); loginErr != nil {
			return account, "", loginErr
		}

		return account, token, nil
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return account, "", err
	}

	if email != "" {
		err = s.Pool.QueryRow(ctx, `
select id, display_name, email, avatar_url, created_at, updated_at, last_login_at
from auth_accounts
where lower(email) = lower($1)
and email <> ''
`, email).Scan(
			&account.ID,
			&account.DisplayName,
			&account.Email,
			&account.AvatarURL,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.LastLoginAt,
		)
		if err == nil {
			_, err = s.Pool.Exec(ctx, `
insert into auth_identities (account_id, provider, provider_user_id, email, display_name, avatar_url, last_login_at)
values ($1, $2, $3, $4, $5, $6, now())
on conflict (provider, provider_user_id) do nothing
`, account.ID, provider, providerUserID, email, displayName, avatarURL)
			if err != nil {
				return account, "", err
			}

			token, sessionErr := s.CreateAuthSession(ctx, account.ID)
			if sessionErr != nil {
				return account, "", sessionErr
			}

			if loginErr := s.touchAuthLogin(ctx, account.ID, provider, providerUserID); loginErr != nil {
				return account, "", loginErr
			}

			return account, token, nil
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return account, "", err
		}
	}

	if displayName == "" {
		displayName = provider + " player"
	}

	displayName, err = s.uniqueAuthDisplayName(ctx, displayName)
	if err != nil {
		return account, "", err
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return account, "", err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
insert into auth_accounts (display_name, email, avatar_url, last_login_at)
values ($1, $2, $3, now())
returning id, display_name, email, avatar_url, created_at, updated_at, last_login_at
`, displayName, email, avatarURL).Scan(
		&account.ID,
		&account.DisplayName,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
	)
	if err != nil {
		return account, "", err
	}

	_, err = tx.Exec(ctx, `
insert into auth_identities (account_id, provider, provider_user_id, email, display_name, avatar_url, last_login_at)
values ($1, $2, $3, $4, $5, $6, now())
`, account.ID, provider, providerUserID, email, displayName, avatarURL)
	if err != nil {
		return account, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return account, "", err
	}

	token, err := s.CreateAuthSession(ctx, account.ID)
	if err != nil {
		return account, "", err
	}

	return account, token, nil
}

func (s *Store) CreateAuthSession(ctx context.Context, accountID int64) (string, error) {
	token, err := randomAuthToken()
	if err != nil {
		return "", err
	}

	_, err = s.Pool.Exec(ctx, `
insert into auth_sessions (session_hash, account_id, expires_at)
values ($1, $2, now() + $3::interval)
`, hashAuthToken(token), accountID, authSessionDuration.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Store) AccountByAuthSession(ctx context.Context, token string) (AuthAccount, error) {
	var account AuthAccount

	tokenHash := hashAuthToken(strings.TrimSpace(token))
	if tokenHash == "" {
		return account, ErrAuthInvalidSession
	}

	err := s.Pool.QueryRow(ctx, `
select a.id, a.display_name, a.email, a.avatar_url, a.created_at, a.updated_at, a.last_login_at
from auth_sessions s
join auth_accounts a on a.id = s.account_id
where s.session_hash = $1
and s.expires_at > now()
`, tokenHash).Scan(
		&account.ID,
		&account.DisplayName,
		&account.Email,
		&account.AvatarURL,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return account, ErrAuthInvalidSession
	}
	if err != nil {
		return account, err
	}

	_, _ = s.Pool.Exec(ctx, `
update auth_sessions
set last_seen_at = now()
where session_hash = $1
`, tokenHash)

	return account, nil
}

func (s *Store) DeleteAuthSession(ctx context.Context, token string) error {
	_, err := s.Pool.Exec(ctx, `
delete from auth_sessions
where session_hash = $1
`, hashAuthToken(token))
	return err
}

func (s *Store) CreateOAuthState(ctx context.Context, provider, redirectPath string) (string, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	redirectPath = strings.TrimSpace(redirectPath)

	if redirectPath == "" || !strings.HasPrefix(redirectPath, "/") || strings.HasPrefix(redirectPath, "//") {
		redirectPath = "/"
	}

	state, err := randomAuthToken()
	if err != nil {
		return "", err
	}

	_, err = s.Pool.Exec(ctx, `
insert into auth_oauth_states (state_hash, provider, redirect_path, expires_at)
values ($1, $2, $3, now() + $4::interval)
`, hashAuthToken(state), provider, redirectPath, authStateDuration.String())
	if err != nil {
		return "", err
	}

	return state, nil
}

func (s *Store) ConsumeOAuthState(ctx context.Context, provider, state string) (string, error) {
	var redirectPath string

	err := s.Pool.QueryRow(ctx, `
delete from auth_oauth_states
where state_hash = $1
and provider = $2
and expires_at > now()
returning redirect_path
`, hashAuthToken(state), strings.TrimSpace(strings.ToLower(provider))).Scan(&redirectPath)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrAuthInvalidState
	}
	if err != nil {
		return "", err
	}

	_, _ = s.Pool.Exec(ctx, `
delete from auth_oauth_states
where expires_at <= now()
`)

	return redirectPath, nil
}

func (s *Store) touchAuthLogin(ctx context.Context, accountID int64, provider, providerUserID string) error {
	_, err := s.Pool.Exec(ctx, `
update auth_accounts
set last_login_at = now(),
updated_at = now()
where id = $1
`, accountID)
	if err != nil {
		return err
	}

	_, err = s.Pool.Exec(ctx, `
update auth_identities
set last_login_at = now()
where account_id = $1
and provider = $2
and provider_user_id = $3
`, accountID, provider, providerUserID)

	return err
}

func (s *Store) uniqueAuthDisplayName(ctx context.Context, requested string) (string, error) {
	base := cleanAuthDisplayName(requested)
	if base == "" || IsReservedAuthDisplayName(base) {
		base = "NamiFan"
	}

	for i := 0; i < 100; i++ {
		candidate := base
		if i > 0 {
			candidate = fmt.Sprintf("%s%d", base, i+1)
		}

		if IsReservedAuthDisplayName(candidate) {
			continue
		}

		var existingID int64
		err := s.Pool.QueryRow(ctx, `
select id
from auth_accounts
where lower(display_name) = lower($1)
`, candidate).Scan(&existingID)
		if errors.Is(err, pgx.ErrNoRows) {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
	}

	suffix, err := randomAuthToken()
	if err != nil {
		return "", err
	}

	if len(suffix) > 8 {
		suffix = suffix[:8]
	}

	return fmt.Sprintf("%s%s", base, suffix), nil
}

func randomAuthToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashAuthToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}

	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func normalizeAuthEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func cleanAuthDisplayName(displayName string) string {
	displayName = strings.TrimSpace(displayName)
	displayName = strings.Join(strings.Fields(displayName), " ")

	if len(displayName) > 32 {
		displayName = displayName[:32]
	}

	return displayName
}

func (s *Store) AuthAccountHasVerifiedProviderEmail(ctx context.Context, accountID int64, provider, email string) (bool, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	email = normalizeAuthEmail(email)

	if accountID < 1 || provider == "" || email == "" {
		return false, nil
	}

	var allowed bool
	err := s.Pool.QueryRow(ctx, `
select exists (
select 1
from auth_identities
where account_id = $1
and provider = $2
and lower(email) = lower($3)
)
`, accountID, provider, email).Scan(&allowed)
	if err != nil {
		return false, err
	}

	return allowed, nil
}

func IsReservedAuthDisplayName(displayName string) bool {
	normalized := reservedDisplayNameKey(displayName)

	if normalized == "" {
		return false
	}

	switch normalized {
	case "soryn", "nami", "namichan":
		return true
	default:
		return false
	}
}

func reservedDisplayNameKey(displayName string) string {
	displayName = strings.TrimSpace(strings.ToLower(displayName))

	replacer := strings.NewReplacer(
		"0", "o",
		"1", "i",
		"3", "e",
		"4", "a",
		"@", "a",
		"5", "s",
		"$", "s",
		"7", "t",
	)

	displayName = replacer.Replace(displayName)

	var builder strings.Builder
	for _, r := range displayName {
		if r >= 'a' && r <= 'z' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
