BEGIN;
--- email table primary key uuids should end in 2 digits, 11, 12, 13, etc
--- employer table primary key uuids should end in 3 digits, 201, 202, 203, etc
--- domain table primary key uuids should end in 4 digits, 3001, 3002, 3003, etc
--- org_users table primary key uuids should end in 5 digits, 40001, 40002, 40003, etc
INSERT INTO public.emails (email_key, email_from, email_to, email_cc, email_bcc, email_subject, email_html_body, email_text_body, email_state, created_at, processed_at)
    VALUES ('12345678-0003-0003-0003-000000000011'::uuid, 'no-reply@vetchi.org', ARRAY['admin@location.example'], NULL, NULL, 'Welcome to Vetchi Subject', 'Welcome to Vetchi HTML Body', 'Welcome to Vetchi Text Body', 'PROCESSED', timezone('UTC'::text, now()), timezone('UTC'::text, now()));
INSERT INTO public.employers (id, client_id_type, employer_state, onboard_admin_email, onboard_secret_token, token_valid_till, onboard_email_id, created_at)
    VALUES ('12345678-0003-0003-0003-000000000201'::uuid, 'DOMAIN', 'ONBOARDED', 'admin@location.example', 'blah', timezone('UTC'::text, now()) + interval '1 day', '12345678-0003-0003-0003-000000000011'::uuid, timezone('UTC'::text, now()));
INSERT INTO public.domains (id, domain_name, domain_state, employer_id, created_at)
    VALUES ('12345678-0003-0003-0003-000000003001'::uuid, 'location.example', 'VERIFIED', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now()));
INSERT INTO public.org_users (id, email, password_hash, org_user_roles, org_user_state, employer_id, created_at)
    VALUES ('12345678-0003-0003-0003-000000040001'::uuid, 'admin@location.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['ADMIN']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0003-0003-0003-000000040002'::uuid, 'crud1@location.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['LOCATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0003-0003-0003-000000040003'::uuid, 'crud2@location.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['LOCATIONS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0003-0003-0003-000000040004'::uuid, 'viewer@location.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['LOCATIONS_VIEWER']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0003-0003-0003-000000040005'::uuid, 'non-location@location.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['COST_CENTERS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now())),
    ('12345678-0003-0003-0003-000000040006'::uuid, 'multiple-non-location-roles@location.example', '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', ARRAY['COST_CENTERS_CRUD', 'OPENINGS_CRUD']::org_user_roles[], 'ACTIVE_ORG_USER', '12345678-0003-0003-0003-000000000201'::uuid, timezone('UTC'::text, now()));
COMMIT;