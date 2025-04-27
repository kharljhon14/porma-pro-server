CREATE TABLE work_experiences(
    id bigserial PRIMARY KEY,
    account_id bigint NOT NULL,
    role varchar(255) NOT NULL,
    company varchar(255) NOT NULL,
    location varchar(255) NOT NULL,
    summary varchar(6000) NOT NULL,
    start_date timestamp NOT NULL,
    end_date timestamp
);

ALTER TABLE "work_experiences" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id"); 

CREATE INDEX ON "work_experiences" ("id", "account_id");