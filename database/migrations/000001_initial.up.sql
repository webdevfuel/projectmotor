CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "name" text,
    "email" text NOT NULL,
    "gh_access_token" text NOT NULL,
    "gh_user_id" integer NOT NULL
);

--> statement-breakpoint
CREATE UNIQUE INDEX users_email_idx ON users (email);

--> statement-breakpoint
CREATE TABLE "sessions" (
    "id" serial PRIMARY KEY,
    "user_id" integer NOT NULL,
    "token" text NOT NULL,
    "user_agent" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

