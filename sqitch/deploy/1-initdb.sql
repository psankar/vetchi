-- Deploy vetchi:1-initdb to pg

BEGIN;

CREATE TABLE IF NOT EXISTS hub_users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE email_states AS ENUM ('PENDING', 'PROCESSED');
CREATE TABLE emails(
	id BIGSERIAL PRIMARY KEY,
	email_from TEXT NOT NULL,
	email_to TEXT ARRAY NOT NULL,
	email_cc TEXT ARRAY,
	email_bcc TEXT ARRAY,
	email_subject TEXT NOT NULL,
	email_html_body TEXT NOT NULL,
	email_text_body TEXT NOT NULL,
	email_state email_states NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
	processed_at TIMESTAMP WITH TIME ZONE
);

---

CREATE TYPE client_id_types AS ENUM ('DOMAIN');
CREATE TYPE employer_states AS ENUM (
    'ONBOARD_PENDING',
    'ONBOARDED',
    'DEBOARDED'
);
CREATE TABLE employers (
    id BIGSERIAL PRIMARY KEY,
    client_id_type client_id_types NOT NULL,
    employer_state employer_states NOT NULL,
    onboard_admin_email TEXT NOT NULL,

    -- TODO: Perhaps we can move this to org_user_tokens ?
    onboard_secret_token TEXT,
    token_valid_till TIMESTAMP WITH TIME ZONE,

    --- Despite its name, it should not be confused with an email address. 
    --- This is the rowid in the 'emails' table for the welcome email sent.
    onboard_email_id BIGINT REFERENCES emails(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE domain_states AS ENUM (
    'VERIFIED',
    'DEBOARDED'
);
CREATE TABLE domains (
    id BIGSERIAL PRIMARY KEY,
    domain_name TEXT NOT NULL UNIQUE,
    domain_state domain_states NOT NULL,

    employer_id INTEGER REFERENCES employers(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE org_user_roles AS ENUM ('ADMIN', 'RECRUITER', 'INTERVIEWER');
CREATE TYPE org_user_states AS ENUM ('ACTIVE', 'LOCKED');
CREATE TABLE org_users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    org_user_role org_user_roles NOT NULL,
    org_user_state org_user_states NOT NULL,
    
    employer_id INTEGER REFERENCES employers(id) NOT NULL,
    UNIQUE (email, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE token_types AS ENUM ('USER_SESSION', 'TGT', 'EMAIL');
CREATE TABLE org_user_tokens (
    token TEXT PRIMARY KEY,
    org_user_id INTEGER REFERENCES org_users(id) NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    token_type token_types NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

COMMIT;
