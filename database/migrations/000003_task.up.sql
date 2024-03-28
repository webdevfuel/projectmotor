CREATE TABLE task (
    "id" serial PRIMARY KEY,
    "title" text NOT NULL,
    "description" text,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    "project_id" integer NOT NULL,
    "owner_id" integer NOT NULL,
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects (id),
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users (id)
);

