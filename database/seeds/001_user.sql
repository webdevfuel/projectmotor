INSERT INTO users ("id", "name", "email", "avatar")
    VALUES (1, 'Web Dev Fuel', 'hello@webdevfuel.com', '');

--> statement-breakpoint
INSERT INTO accounts ("id", "user_id", "access_token")
    VALUES (1, 1, 'REDACTED');

--> statement-breakpoint
INSERT INTO users ("id", "name", "email", "avatar")
    VALUES (2, 'John Doe', 'johndoe@gmail.com', '');

--> statement-breakpoint
INSERT INTO accounts ("id", "user_id", "access_token")
    VALUES (2, 2, 'REDACTED');

