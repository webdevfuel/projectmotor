CREATE TABLE tasks (
    "id" serial PRIMARY KEY,
    "title" text NOT NULL,
    "description" text,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    "project_id" integer,
    "owner_id" integer NOT NULL,
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE SET NULL,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users (id)
);

