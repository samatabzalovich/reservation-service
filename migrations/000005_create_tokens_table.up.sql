Drop type if exists TOKENSCOPE;
CREATE TYPE TOKENSCOPE AS ENUM ('authentication', 'employee_registration');
CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    institution_id bigint NULL REFERENCES institution(id) ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);
