CREATE TABLE job
(
    id         SERIAL PRIMARY KEY,
    code       VARCHAR UNIQUE NOT NULL,
    schedule   VARCHAR        NOT NULL,
    is_active  BOOLEAN        NOT NULL DEFAULT false,
    is_running BOOLEAN        NOT NULL DEFAULT false
);

CREATE TABLE job_history
(
    id             SERIAL PRIMARY KEY,
    job_code       VARCHAR REFERENCES job (code),
    started        TIMESTAMP,
    finished       TIMESTAMP,
    result         VARCHAR,
    result_message VARCHAR
);
