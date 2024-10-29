package hermione

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hermione/costcenter"
	"github.com/psankar/vetchi/api/internal/hermione/locations"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Config struct {
	Employer struct {
		TGTLife          string `json:"tgt_life" validate:"required,min=1"`
		SessionTokLife   string `json:"session_tok_life" validate:"required,min=1"`
		LTSessionTokLife string `json:"lt_session_tok_life" validate:"required,min=1"`
	} `json:"employer" validate:"required"`

	Postgres struct {
		Host string `json:"host" validate:"required,min=1"`
		Port string `json:"port" validate:"required,min=1"`
		User string `json:"user" validate:"required,min=1"`
		DB   string `json:"db" validate:"required,min=1"`
	} `json:"postgres" validate:"required"`

	Port string `json:"port" validate:"required,min=1,number"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("/etc/hermione-config/config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Whatever can be validated by the struct tags, is done here. More
	// validations continue to happen in the New() function
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// We can store the token lifetimes as float64 instead of time.Duration and
// may be able to save some time avoiding the Mins() call everytime we need to
// refer to these values. But a pretty-printed time.Duration is easier to debug
// than a random float64 literal.
type employer struct {
	// Token Granting Token life - Should be used for /employer/signin
	tgtLife time.Duration

	// User Session Tokens life - Should be used for /employer/tfa
	sessionTokLife         time.Duration
	longTermSessionTokLife time.Duration
}

type Hermione struct {
	// These are initialized from configmap
	employer employer
	port     string

	// These are initialized programmatically in New()
	db    db.DB
	log   *slog.Logger
	mw    *middleware.Middleware
	vator *vetchi.Vator
}

func NewHermione() (*Hermione, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD not set")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	pgConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.DB,
		pgPassword,
	)
	db, err := postgres.New(pgConnStr, logger)
	if err != nil {
		return nil, err
	}

	vator, err := vetchi.InitValidator(logger)
	if err != nil {
		return nil, err
	}

	tgtLife, err := time.ParseDuration(config.Employer.TGTLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.TGTLife: %w", err)
	}
	sessionTokLife, err := time.ParseDuration(config.Employer.SessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.SessionTokLife: %w", err)
	}
	ltsTokLife, err := time.ParseDuration(config.Employer.LTSessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("config.Employer.LTSessionTokLife: %w", err)
	}

	return &Hermione{
		db:    db,
		port:  fmt.Sprintf(":%s", config.Port),
		log:   logger,
		mw:    middleware.NewMiddleware(db, logger),
		vator: vator,
		employer: employer{
			tgtLife:                tgtLife,
			sessionTokLife:         sessionTokLife,
			longTermSessionTokLife: ltsTokLife,
		},
	}, nil
}

func (h *Hermione) Run() error {
	// Authentication related endpoints
	http.HandleFunc("/employer/get-onboard-status", h.getOnboardStatus)
	http.HandleFunc("/employer/set-onboard-password", h.setOnboardPassword)
	http.HandleFunc("/employer/signin", h.employerSignin)
	http.HandleFunc("/employer/tfa", h.employerTFA)

	// CostCenter related endpoints
	h.mw.Protect(
		"/employer/add-cost-center",
		http.HandlerFunc(costcenter.AddCostCenter(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-centers",
		http.HandlerFunc(costcenter.GetCostCenters(h)),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.CostCentersCRUD,
			vetchi.CostCentersViewer,
		},
	)
	h.mw.Protect(
		"/employer/defunct-cost-center",
		http.HandlerFunc(costcenter.DefunctCostCenter(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/rename-cost-center",
		http.HandlerFunc(costcenter.RenameCostCenter(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/update-cost-center",
		http.HandlerFunc(costcenter.UpdateCostCenter(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-center",
		http.HandlerFunc(costcenter.GetCostCenter(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersViewer},
	)

	// Location related endpoints
	h.mw.Protect(
		"/employer/add-location",
		http.HandlerFunc(locations.AddLocation(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/defunct-location",
		http.HandlerFunc(locations.DefunctLocation(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/get-locations",
		http.HandlerFunc(locations.GetLocations(h)),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.LocationsCRUD,
			vetchi.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/get-location",
		http.HandlerFunc(locations.GetLocation(h)),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.LocationsCRUD,
			vetchi.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/rename-location",
		http.HandlerFunc(locations.RenameLocation(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/update-location",
		http.HandlerFunc(locations.UpdateLocation(h)),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)

	return http.ListenAndServe(h.port, nil)
}

func (h *Hermione) DB() db.DB {
	return h.db
}

func (h *Hermione) Log() *slog.Logger {
	return h.log
}

func (h *Hermione) Vator() *vetchi.Vator {
	return h.vator
}
