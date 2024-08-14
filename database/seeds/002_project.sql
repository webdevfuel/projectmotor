INSERT INTO projects (title, description, published, owner_id)
    VALUES ('Project 1', '', FALSE, 1);

--> statement-breakpoint
INSERT INTO projects (title, description, published, owner_id)
    VALUES ('Project 2', '', TRUE, 1);

--> statement-breakpoint
INSERT INTO projects (title, description, published, owner_id)
    VALUES ('Project 3', '', FALSE, 2);

--> statement-breakpoint
INSERT INTO projects (title, description, published, owner_id)
    VALUES ('Project 4', '', TRUE, 2);

