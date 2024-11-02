package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type ClientIDType string

const (
	DomainClientIDType ClientIDType = "DOMAIN"
)

type EmployerState string

const (
	OnboardPendingEmployerState EmployerState = "ONBOARD_PENDING"
	OnboardedEmployerState      EmployerState = "ONBOARDED"
	DeboardedEmployerState      EmployerState = "DEBOARDED"
)

type DomainState string

const (
	VerifiedDomainState  DomainState = "VERIFIED"
	DeboardedDomainState DomainState = "DEBOARDED"
)

type Employer struct {
	ID                 uuid.UUID          `db:"id"`
	ClientIDType       ClientIDType       `db:"client_id_type"`
	EmployerState      EmployerState      `db:"employer_state"`
	OnboardAdminEmail  string             `db:"onboard_admin_email"`
	OnboardSecretToken pgtype.Text        `db:"onboard_secret_token"`
	TokenValidTill     pgtype.Timestamptz `db:"token_valid_till"`
	OnboardEmailID     uuid.UUID          `db:"onboard_email_id"`
	CreatedAt          time.Time          `db:"created_at"`
}

type Domain struct {
	ID          uuid.UUID   `db:"id"`
	DomainName  string      `db:"domain_name"`
	DomainState DomainState `db:"domain_state"`
	EmployerID  uuid.UUID   `db:"employer_id"`
	CreatedAt   time.Time   `db:"created_at"`
}

type OnboardInfo struct {
	EmployerID     uuid.UUID
	AdminEmailAddr string
	DomainName     string
}

type OnboardReq struct {
	DomainName string
	Password   string
	Token      string
}

type OrgUserAuth struct {
	OrgUserID     uuid.UUID
	OrgUserEmail  string
	EmployerID    uuid.UUID
	OrgUserRoles  []vetchi.OrgUserRole
	PasswordHash  string
	EmployerState EmployerState
	OrgUserState  vetchi.OrgUserState
}

type OrgUserCreds struct {
	ClientID string
	Email    string
}

type EmployerTFA struct {
	TFACode  TokenReq
	TFAToken TokenReq
	Email    Email
}
