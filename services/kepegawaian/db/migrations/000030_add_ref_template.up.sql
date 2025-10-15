BEGIN;

CREATE TABLE "ref_template" (
	id serial PRIMARY KEY,
	nama varchar(255) NOT NULL,
	file_base64 text,
	created_at timestamptz DEFAULT now(),
	updated_at timestamptz DEFAULT now(),
	deleted_at timestamptz
);

COMMIT;
