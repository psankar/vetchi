BEGIN;

DELETE FROM org_user_tokens WHERE org_user_id IN (SELECT id FROM org_users WHERE employer_id = '12345678-0002-0002-0002-000000000201'::UUID);
DELETE FROM public.org_users WHERE employer_id = '12345678-0002-0002-0002-000000000201'::UUID;
DELETE FROM public.org_cost_centers WHERE employer_id = '12345678-0002-0002-0002-000000000201'::UUID;
DELETE FROM public.domains WHERE employer_id = '12345678-0002-0002-0002-000000000201'::UUID;
DELETE FROM public.employers WHERE id = '12345678-0002-0002-0002-000000000201'::UUID;
DELETE FROM public.emails WHERE email_key = '12345678-0002-0002-0002-000000000011'::UUID;

COMMIT;