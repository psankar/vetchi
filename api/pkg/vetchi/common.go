package vetchi

type CountryCode string
type Currency string
type EmailAddress string
type Password string
type TimeZone string

type ValidationErrors struct {
	Errors []string `json:"errors"`
}

type EmployerSignInResponse struct {
	Token string `json:"token"`
}
