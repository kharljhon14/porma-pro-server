CREATE TABLE IF NOT EXISTS personal_infos(
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "full_name" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "phone_number" varchar(255) NOT NULL,
    "linkedin_url" varchar(255),
    "personal_url" varchar(255),
    "country" varchar(255) NOT NULL,
    "state" varchar(255) NOT NULL,
    "city" varchar(255) NOT NULL
);

ALTER TABLE "personal_infos" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "personal_infos" ("id", "account_id");