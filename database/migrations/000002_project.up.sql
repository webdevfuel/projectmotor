CREATE TABLE "projects" (
    "id" serial PRIMARY KEY,
    "title" text NOT NULL,
    "description" text,
    "published" boolean NOT NULL DEFAULT FALSE,
    "owner_id" integer NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users (id)
);

