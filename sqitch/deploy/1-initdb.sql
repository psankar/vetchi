-- Deploy vetchi:1-initdb to pg

ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_min_duration_statement = 0;
ALTER SYSTEM SET log_duration = 'on';
SELECT pg_reload_conf();

BEGIN;

CREATE TYPE hub_user_states AS ENUM ('ACTIVE_HUB_USER', 'DISABLED_HUB_USER', 'DELETED_HUB_USER');
CREATE TABLE IF NOT EXISTS hub_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    handle TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    state hub_user_states NOT NULL,
    resident_country_code TEXT NOT NULL,
    resident_city TEXT,
    preferred_language TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE hub_user_token_types AS ENUM (
    -- Sent as response to the TFA API
    'HUB_USER_SESSION',
    'HUB_USER_LTS',

    -- Sent as response to the Login API
    'HUB_USER_TFA_TOKEN',

    -- Sent as response to the Reset Password API
    'HUB_USER_RESET_PASSWORD_TOKEN'
);

CREATE TABLE hub_user_tokens (
    token TEXT CONSTRAINT hub_user_tokens_pkey PRIMARY KEY,
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,
    token_type hub_user_token_types NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE hub_user_tfa_codes (
    code TEXT NOT NULL,
    tfa_token TEXT NOT NULL REFERENCES hub_user_tokens(token) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE official_email_states AS ENUM ('PENDING', 'VERIFIED');
CREATE TABLE hub_users_official_emails (
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,
    official_email TEXT NOT NULL,
    state official_email_states NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    PRIMARY KEY (hub_user_id, official_email)
);

CREATE TYPE email_states AS ENUM ('PENDING', 'PROCESSED');
CREATE TABLE emails(
	email_key UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id_type client_id_types NOT NULL,
    employer_state employer_states NOT NULL,
    company_name TEXT NOT NULL,

    onboard_admin_email TEXT NOT NULL,

    -- TODO: Perhaps we can move this to org_user_tokens ?
    onboard_secret_token TEXT,
    token_valid_till TIMESTAMP WITH TIME ZONE,

    --- Despite its name, it should not be confused with an email address.
    --- This is the rowid in the 'emails' table for the welcome email sent.
    onboard_email_id UUID REFERENCES emails(email_key),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE domain_states AS ENUM (
    'VERIFIED',
    'DEBOARDED'
);
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    domain_name TEXT NOT NULL,
    CONSTRAINT uniq_domain_name UNIQUE (domain_name),

    domain_state domain_states NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    CONSTRAINT uniq_employer_domain_id UNIQUE (employer_id, id)
);

CREATE TABLE employer_primary_domains(
    employer_id UUID NOT NULL REFERENCES employers(id) ON DELETE CASCADE,
    domain_id UUID NOT NULL REFERENCES domains(id) ON DELETE CASCADE,

    PRIMARY KEY (employer_id),
    CONSTRAINT fk_employer_domain_match FOREIGN KEY (employer_id, domain_id) 
        REFERENCES domains(employer_id, id)
);
---

CREATE TYPE org_user_roles AS ENUM (
    'ADMIN',
    'COST_CENTERS_CRUD',
    'COST_CENTERS_VIEWER',
    'LOCATIONS_CRUD',
    'LOCATIONS_VIEWER',
    'ORG_USERS_CRUD',
    'ORG_USERS_VIEWER',
    'OPENINGS_CRUD',
    'OPENINGS_VIEWER'
);
CREATE TYPE org_user_states AS ENUM (
    'ACTIVE_ORG_USER',
    'INVITED_ORG_USER',
    'ADDED_ORG_USER',
    'DISABLED_ORG_USER',
    'REPLICATED_ORG_USER'
);
CREATE TABLE org_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT,
    org_user_roles org_user_roles[] NOT NULL,
    org_user_state org_user_states NOT NULL,

--- As of now, we have only one org per employer. This may change in future.
    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_email_employer_id UNIQUE (email, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE org_user_token_types AS ENUM (
    -- Sent as response to the TFA API
    'EMPLOYER_SESSION',
    'EMPLOYER_LTS',

    -- Sent as response to the SignIn API
    'EMPLOYER_TFA_TOKEN'
);
CREATE TABLE org_user_tokens (
    token TEXT CONSTRAINT org_user_tokens_pkey PRIMARY KEY,
    org_user_id UUID REFERENCES org_users(id) NOT NULL,
    token_type org_user_token_types NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE org_user_tfa_codes (
    code TEXT NOT NULL,
    tfa_token TEXT NOT NULL REFERENCES org_user_tokens(token) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TABLE org_user_invites (
    token TEXT CONSTRAINT org_user_invites_pkey PRIMARY KEY,
    org_user_id UUID REFERENCES org_users(id) NOT NULL,
    token_valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE cost_center_states AS ENUM ('ACTIVE_CC', 'DEFUNCT_CC');
CREATE TABLE org_cost_centers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cost_center_name TEXT NOT NULL,
    cost_center_state cost_center_states NOT NULL,
    notes TEXT NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_cost_center_name_employer_id UNIQUE (cost_center_name, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE location_states AS ENUM ('ACTIVE_LOCATION', 'DEFUNCT_LOCATION');
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    country_code TEXT NOT NULL,
    postal_address TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    openstreetmap_url TEXT,
    city_aka TEXT ARRAY,

    location_state location_states NOT NULL,

    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT uniq_location_title_employer_id UNIQUE (title, employer_id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

---

CREATE TYPE opening_states AS ENUM ('DRAFT_OPENING_STATE', 'ACTIVE_OPENING_STATE', 'SUSPENDED_OPENING_STATE', 'CLOSED_OPENING_STATE');
CREATE TYPE opening_types AS ENUM ('FULL_TIME_OPENING', 'PART_TIME_OPENING', 'CONTRACT_OPENING', 'INTERNSHIP_OPENING', 'UNSPECIFIED_OPENING');
CREATE TYPE education_levels AS ENUM ('BACHELOR_EDUCATION', 'MASTER_EDUCATION', 'DOCTORATE_EDUCATION', 'NOT_MATTERS_EDUCATION', 'UNSPECIFIED_EDUCATION');
CREATE TABLE openings (
    employer_id UUID REFERENCES employers(id) NOT NULL,
    id TEXT NOT NULL,
    CONSTRAINT openings_pk PRIMARY KEY (employer_id, id),
    CONSTRAINT opening_id_format_check CHECK (id ~ '^[0-9]{4}-[A-Z][a-z]{2}-[0-9]{2}-[0-9]+$'),

    title TEXT NOT NULL,
    positions INTEGER NOT NULL,
    jd TEXT NOT NULL,
    recruiter UUID REFERENCES org_users(id) NOT NULL,
    hiring_manager UUID REFERENCES org_users(id) NOT NULL,
    cost_center_id UUID REFERENCES org_cost_centers(id),
    employer_notes TEXT,
    remote_country_codes TEXT[],
    remote_timezones TEXT[],
    opening_type opening_types NOT NULL,
    yoe_min INTEGER NOT NULL,
    yoe_max INTEGER NOT NULL,
    min_education_level education_levels NOT NULL,
    salary_min NUMERIC,
    salary_max NUMERIC,
    salary_currency TEXT,
    state opening_states NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),
    last_updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now()),

    pagination_key BIGSERIAL
);

CREATE TABLE opening_hiring_team(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    hiring_team_mate_id UUID REFERENCES org_users(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, hiring_team_mate_id)
);

CREATE TABLE opening_locations(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    location_id UUID REFERENCES locations(id) NOT NULL,

    PRIMARY KEY (employer_id, opening_id, location_id)
);

CREATE TABLE opening_watchers(
    employer_id UUID NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    watcher_id UUID REFERENCES org_users(id) NOT NULL,
    PRIMARY KEY (employer_id, opening_id, watcher_id)
);

CREATE TYPE application_color_tags AS ENUM ('GREEN', 'YELLOW', 'RED');
CREATE TYPE application_states AS ENUM ('APPLIED', 'REJECTED', 'SHORTLISTED', 'WITHDRAWN', 'EXPIRED');
CREATE TABLE applications (
    id TEXT PRIMARY KEY,
    employer_id UUID REFERENCES employers(id) NOT NULL,
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),
    cover_letter TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    internal_filename TEXT NOT NULL,
    application_state application_states NOT NULL,

    color_tag application_color_tags,

    -- The user who applied for the opening
    hub_user_id UUID REFERENCES hub_users(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

CREATE TYPE candidacy_states AS ENUM (
    'INTERVIEWING', 
    'OFFERED', 'OFFER_DECLINED', 'OFFER_ACCEPTED', 
    'CANDIDATE_UNSUITABLE',
    'CANDIDATE_NOT_RESPONDING',
    'EMPLOYER_DEFUNCT'
);
CREATE TABLE candidacies(
    id TEXT PRIMARY KEY,
    application_id TEXT REFERENCES applications(id) NOT NULL,
    CONSTRAINT fk_application FOREIGN KEY (application_id) REFERENCES applications(id),
    employer_id UUID REFERENCES employers(id) NOT NULL,
    CONSTRAINT fk_employer FOREIGN KEY (employer_id) REFERENCES employers(id),
    CONSTRAINT uniq_application_employer_id UNIQUE (application_id, employer_id),
    opening_id TEXT NOT NULL,
    CONSTRAINT fk_opening FOREIGN KEY (employer_id, opening_id) REFERENCES openings (employer_id, id),

    candidacy_state candidacy_states NOT NULL,

    created_by UUID REFERENCES org_users(id) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('UTC', now())
);

COMMIT;
