Drop type if exists USERTYPE;
Create TYPE USERTYPE as ENUM ('client', 'owner', 'admin', 'employee');
CREATE TABLE IF NOT EXISTS users (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
userName text NOT NULL,
type USERTYPE NOT NULL,
number varchar(12) UNIQUE NOT NULL,
password_hash bytea NOT NULL,
activated bool NOT NULL,
version integer NOT NULL DEFAULT 1
);