CREATE TABLE "accounts"(
    "id" bigserial PRIMARY KEY,
    "email" varchar UNIQUE NOT NULL,
    "password_hash" varchar,
    "sso_provider" varchar,
    "sso_provider_id" varchar,
    "full_name" varchar NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now()),
    "is_verified" boolean NOT NULL DEFAULT False,
);