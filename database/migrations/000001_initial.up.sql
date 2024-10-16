CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "name" text,
    "email" text NOT NULL,
    "gh_access_token" text NOT NULL,
    "gh_user_id" integer NOT NULL
);

--> statement-breakpoint
CREATE TABLE "sessions" (
    "id" serial PRIMARY KEY,
    "token" text NOT NULL,
    "user_id" integer NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

