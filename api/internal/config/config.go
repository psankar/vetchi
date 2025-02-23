package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type HermioneConfigOnDisk struct {
	Employer struct {
		TFATokLife     string `json:"tfa_tok_life" validate:"required"`
		SessionTokLife string `json:"session_tok_life" validate:"required"`
		LTSTokLife     string `json:"lts_tok_life" validate:"required"`
		InviteTokLife  string `json:"employee_invite_tok_life" validate:"required"`
	} `json:"employer" validate:"required"`

	Hub struct {
		WebURL               string `json:"web_url" validate:"required"`
		TFATokLife           string `json:"tfa_tok_life" validate:"required"`
		SessionTokLife       string `json:"session_tok_life" validate:"required"`
		LTSTokLife           string `json:"lts_tok_life" validate:"required"`
		InviteTokLife        string `json:"hub_user_invite_tok_life" validate:"required"`
		PasswordResetTokLife string `json:"password_reset_tok_life" validate:"required"`
	} `json:"hub" validate:"required"`

	S3 struct {
		Endpoint  string `json:"endpoint" validate:"required"`
		Region    string `json:"region" validate:"required"`
		Bucket    string `json:"bucket" validate:"required"`
		AccessKey string `json:"access_key" validate:"required"`
	} `json:"s3" validate:"required"`

	Port string `json:"port" validate:"required,min=1,number"`

	TimingAttackDelay string `json:"timing_attack_delay" validate:"required"`
}

type Hermione struct {
	Employer struct {
		TFATokLife     time.Duration
		SessionTokLife time.Duration
		LTSTokLife     time.Duration
		InviteTokLife  time.Duration
	}

	Hub struct {
		WebURL               string
		TFATokLife           time.Duration
		SessionTokLife       time.Duration
		LTSTokLife           time.Duration
		InviteTokLife        time.Duration
		PasswordResetTokLife time.Duration
	}

	S3 struct {
		Endpoint  string
		Region    string
		Bucket    string
		AccessKey string
		SecretKey string
	}

	Port              int
	TimingAttackDelay time.Duration
}

func LoadHermioneConfig() (*Hermione, error) {
	data, err := os.ReadFile("/etc/hermione-config/config.json")
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cmap := &HermioneConfigOnDisk{}
	if err := json.Unmarshal(data, cmap); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Whatever can be validated by the struct tags, is done here. More
	// validations continue to happen below
	validate := validator.New()
	if err := validate.Struct(cmap); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	hc := &Hermione{}

	// Load S3 secret key from environment
	hc.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
	if hc.S3.SecretKey == "" {
		return nil, fmt.Errorf("S3_SECRET_KEY environment variable is required")
	}

	// Copy S3 config
	s3 := cmap.S3
	hc.S3.Endpoint = s3.Endpoint
	hc.S3.Region = s3.Region
	hc.S3.Bucket = s3.Bucket
	hc.S3.AccessKey = s3.AccessKey

	hc.Port, err = strconv.Atoi(cmap.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %w", err)
	}

	hc.TimingAttackDelay, err = time.ParseDuration(cmap.TimingAttackDelay)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timing attack delay: %w", err)
	}

	emp := cmap.Employer
	hc.Employer.TFATokLife, err = time.ParseDuration(emp.TFATokLife)
	if err != nil {
		return nil, fmt.Errorf("employer tfa token life: %w", err)
	}
	hc.Employer.SessionTokLife, err = time.ParseDuration(emp.SessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("employer session token life: %w", err)
	}
	hc.Employer.LTSTokLife, err = time.ParseDuration(emp.LTSTokLife)
	if err != nil {
		return nil, fmt.Errorf("employer lts token life: %w", err)
	}
	hc.Employer.InviteTokLife, err = time.ParseDuration(emp.InviteTokLife)
	if err != nil {
		return nil, fmt.Errorf("employer invite token life: %w", err)
	}

	hub := cmap.Hub
	hc.Hub.WebURL = hub.WebURL
	hc.Hub.TFATokLife, err = time.ParseDuration(hub.TFATokLife)
	if err != nil {
		return nil, fmt.Errorf("hub tfa token life: %w", err)
	}
	hc.Hub.SessionTokLife, err = time.ParseDuration(hub.SessionTokLife)
	if err != nil {
		return nil, fmt.Errorf("hub session token life: %w", err)
	}
	hc.Hub.LTSTokLife, err = time.ParseDuration(hub.LTSTokLife)
	if err != nil {
		return nil, fmt.Errorf("hub lts token life: %w", err)
	}
	hc.Hub.InviteTokLife, err = time.ParseDuration(hub.InviteTokLife)
	if err != nil {
		return nil, fmt.Errorf("hub invite token life: %w", err)
	}
	hc.Hub.PasswordResetTokLife, err = time.ParseDuration(
		hub.PasswordResetTokLife,
	)
	if err != nil {
		return nil, fmt.Errorf("hub password reset token life: %w", err)
	}

	return hc, nil
}
