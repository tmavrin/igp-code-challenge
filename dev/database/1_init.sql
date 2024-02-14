CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION update_modified_column()
	RETURNS trigger LANGUAGE plpgsql AS $function$
BEGIN
    NEW.modified = now();
    RETURN NEW; 
END;
$function$;

CREATE TABLE IF NOT EXISTS accounts (
	id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	email CITEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER accounts_modtime BEFORE UPDATE
	ON accounts 
	FOR EACH ROW EXECUTE PROCEDURE update_modified_column();