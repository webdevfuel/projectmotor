CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "name" text,
    "email" text NOT NULL,
    "avatar" text
);

CREATE TABLE "accounts" (
    "user_id" integer NOT NULL,
    "account_id" text PRIMARY KEY,
    "access_token" text NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

