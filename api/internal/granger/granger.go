package granger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Config struct {
	Env string `json:"env" validate:"required,min=1"`

	OnboardTokenLife string `json:"onboard_token_life" validate:"required,min=1"`

	Port string `json:"port" validate:"required,min=1,number"`

	SMTP struct {
		Host string `json:"host" validate:"required,min=1"`
		Port string `json:"port" validate:"required,min=1,number"`
		User string `json:"user" validate:"required,min=1"`
	} `json:"smtp" validate:"required"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("/etc/granger-config/config.json")
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

	if config.Env != vetchi.ProdEnv && config.Env != vetchi.DevEnv {
		return nil, fmt.Errorf(
			"%q is not one of [%q, %q]",
			config.Env,
			vetchi.ProdEnv,
			vetchi.DevEnv,
		)
	}

	return config, nil
}

type smtp struct {
	host     string
	port     int
	user     string
	password string
}

type Granger struct {
	// These are initialized from configmap
	env              string
	onboardTokenLife time.Duration
	port             string
	smtp             smtp

	// These are initialized programatically in NewGranger()
	db  db.DB
	log util.Logger
	wg  sync.WaitGroup
}

func NewGranger() (*Granger, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	pgConnStr := os.Getenv("POSTGRES_URI")
	if pgConnStr == "" {
		return nil, fmt.Errorf("POSTGRES_URI not set")
	}

	logger := util.Logger{
		Log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})),
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD not set")
	}

	db, err := postgres.New(pgConnStr, logger)
	if err != nil {
		return nil, err
	}

	smtpPort, err := strconv.Atoi(config.SMTP.Port)
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT is invalid: %w", err)
	}

	tokenDuration, err := time.ParseDuration(config.OnboardTokenLife)
	if err != nil {
		return nil, fmt.Errorf("OnboardTokenLife is invalid: %w", err)
	}

	return &Granger{
		env:  config.Env,
		port: fmt.Sprintf(":%s", config.Port),
		smtp: smtp{
			host:     config.SMTP.Host,
			port:     smtpPort,
			user:     config.SMTP.User,
			password: smtpPassword,
		},
		onboardTokenLife: tokenDuration,

		db:  db,
		log: logger,
	}, nil
}

func (g *Granger) Run() error {
	g.wg.Add(1)
	pruneTokensQuit := make(chan struct{})
	go g.pruneTokens(pruneTokensQuit)

	g.wg.Add(1)
	createOnboardEmailsQuit := make(chan struct{})
	go g.createOnboardEmails(createOnboardEmailsQuit)

	g.wg.Add(1)
	pruneOfficialEmailCodesQuit := make(chan struct{})
	go g.pruneOfficialEmailCodes(pruneOfficialEmailCodesQuit)

	g.wg.Add(1)
	mailSenderQuit := make(chan struct{})
	go g.mailSender(mailSenderQuit)

	go func() {
		// For now, we don't have any routes to serve
		// but we will keep this around for future use

		err := http.ListenAndServe(g.port, nil)
		if err != nil {
			g.log.Err("Failed to start HTTP server", "error", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(pruneTokensQuit)
		close(createOnboardEmailsQuit)
		close(mailSenderQuit)
	}()

	g.wg.Wait()
	return nil
}
