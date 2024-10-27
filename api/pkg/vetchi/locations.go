package vetchi

type LocationState string

const (
	ActiveLocation  LocationState = "ACTIVE_LOCATION"
	DefunctLocation LocationState = "DEFUNCT_LOCATION"
)

type Location struct {
	Title            string   `json:"title"`
	CountryCode      string   `json:"country_code"`
	PostalAddress    string   `json:"postal_address"`
	PostalCode       string   `json:"postal_code"`
	OpenStreetMapURL string   `json:"openstreetmap_url"`
	CityAka          []string `json:"city_aka"`

	State LocationState `json:"state"`
}

type AddLocationRequest struct {
	Title            string   `json:"title"             validate:"required,min=3,max=32"`
	CountryCode      string   `json:"country_code"      validate:"required,len=3,validate_country_code"`
	PostalAddress    string   `json:"postal_address"    validate:"required,min=3,max=1024"`
	PostalCode       string   `json:"postal_code"       validate:"required,min=3,max=16"`
	OpenStreetMapURL string   `json:"openstreetmap_url" validate:"omitempty,url,max=255"`
	CityAka          []string `json:"city_aka"          validate:"omitempty,maxItems=3,validate_city_aka"`
}

type DefunctLocationRequest struct {
	Title string `json:"title" validate:"required,min=3,max=32"`
}

type GetLocationRequest struct {
	Title string `json:"title" validate:"required,min=3,max=32"`
}

type GetLocationsRequest struct {
	PaginationKey string `json:"pagination_key"`
	Limit         int    `json:"limit"          validate:"min=0,max=100"`
}

type RenameLocationRequest struct {
	OldTitle string `json:"old_title" validate:"required,min=3,max=32"`
	NewTitle string `json:"new_title" validate:"required,min=3,max=32"`
}

type UpdateLocationRequest struct {
	Title            string   `json:"title"             validate:"required,min=3,max=32"`
	CountryCode      string   `json:"country_code"      validate:"required,len=3,validate_country_code"`
	PostalAddress    string   `json:"postal_address"    validate:"required,min=3,max=1024"`
	PostalCode       string   `json:"postal_code"       validate:"required,min=3,max=16"`
	OpenStreetMapURL string   `json:"openstreetmap_url" validate:"omitempty,url,max=255"`
	CityAka          []string `json:"city_aka"          validate:"omitempty,maxItems=3,validate_city_aka"`
}