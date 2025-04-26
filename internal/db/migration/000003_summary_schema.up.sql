CREATE TABLE summaries (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "summary" varchar(3000) NOT NULL
);

ALTER TABLE "summaries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "summaries" ("id", "account_id");