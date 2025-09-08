BEGIN;

CREATE TABLE "user" (
  id uuid not null,
  source varchar(10) not null,
  nip varchar(20) not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz,
  PRIMARY KEY (id, source)
);

CREATE TABLE role (
  id serial PRIMARY KEY,
  service varchar(50) not null,
  nama varchar(100) not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE TABLE user_role (
  id serial PRIMARY KEY,
  nip varchar(20) not null,
  role_id int4 not null REFERENCES role (id),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  deleted_at timestamptz
);

CREATE UNIQUE INDEX user_role_nip_role_id_unique_idx ON user_role(nip, role_id) WHERE deleted_at IS NULL;

COMMIT;
