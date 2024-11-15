BEGIN;
DELETE FROM opening_hiring_team
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM opening_locations
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM openings
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM locations
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM org_cost_centers
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM org_user_tokens
WHERE org_user_id IN (
    SELECT id FROM org_users 
    WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid
);

DELETE FROM org_users
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM domains
WHERE employer_id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM employers
WHERE id = '12345678-0005-0005-0005-000000000201'::uuid;

DELETE FROM emails
WHERE email_key = '12345678-0005-0005-0005-000000000011'::uuid;

COMMIT;