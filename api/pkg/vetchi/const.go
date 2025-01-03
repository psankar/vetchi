package vetchi

const (
	HubBaseURL      = "https://vetchi.org"
	EmployerBaseURL = "https://employer.vetchi.org"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

const (
	// Sent in the email to the org users
	InviteTokenLenBytes = 16

	// Triggered by the forgot password request and sent to the user's email
	PasswordResetTokenLenBytes = 16

	// Sent as a response to the signin request
	// Used for the /employer/tfa request body
	TGTokenLenBytes = 32

	// Used for the email code that is sent to the user's email for tfa
	EmailTokenLenBytes = 2

	// Used for the session tokens
	SessionTokenLenBytes = 8

	ApplicationIDLenBytes = 16

	CandidacyIDLenBytes = 16

	InterviewIDLenBytes = 16
)

const (
	EmailFrom = "no-reply@vetchi.org"
)
