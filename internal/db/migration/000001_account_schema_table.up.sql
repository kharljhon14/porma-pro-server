CREATE TABLE "accounts"(
    "id" bigserial PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "password_hash" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now()),
    "is_verified" boolean NOT NULL DEFAULT false
);


CREATE INDEX ON "accounts" ("email", "id");