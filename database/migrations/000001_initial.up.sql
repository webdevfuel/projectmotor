CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "name" text,
    "email" text NOT NULL,
    "avatar" text
);

CREATE TABLE "accounts" (
    "id" integer PRIMARY KEY,
    "user_id" integer NOT NULL,
    "access_token" text NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

