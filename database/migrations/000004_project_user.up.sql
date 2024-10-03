CREATE TABLE projects_users (
    project_id integer NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT fk_project FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX projects_users_project_id_user_id_idx ON projects_users (project_id, user_id);

