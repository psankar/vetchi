package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

type HubUserTFA struct {
	TFACode  string
	TFAToken HubTokenReq
	Email    Email
}

type HubUserTO struct {
	ID           uuid.UUID           `db:"id"`
	FullName     string              `db:"full_name"`
	Handle       string              `db:"handle"`
	Email        common.EmailAddress `db:"email"`
	PasswordHash string              `db:"password_hash"`
	State        hub.HubUserState    `db:"state"`
	CreatedAt    time.Time           `db:"created_at"`
	UpdatedAt    time.Time           `db:"updated_at"`
}

type HubUserInitPasswordReset struct {
	Email Email
	HubTokenReq
}

type HubUserPasswordReset struct {
	Token        string
	PasswordHash string
}

type ApplyOpeningReq struct {
	ApplicationID          string
	OpeningIDWithinCompany string
	CompanyDomain          string
	CoverLetter            string
	ResumeSHA              string
	EndorserHandles        []common.Handle
	EndorsementEmails      []Email
}

type OnboardHubUserReq struct {
	// Comes from the user
	InviteToken         string
	FullName            string
	PasswordHash        string
	Tier                hub.HubUserTier
	ResidentCountryCode common.CountryCode
	PreferredLanguage   string
	ShortBio            string
	LongBio             string

	// Generated internally in the handler
	SessionToken                 string
	SessionTokenValidityDuration time.Duration
	SessionTokenType             HubTokenType
}
