INSERT INTO users (email) VALUES ('phat@gmail.com');
INSERT INTO users (email) VALUES ('phat1@gmail.com');
INSERT INTO users (email) VALUES ('phat2@gmail.com');
INSERT INTO users (email) VALUES ('phat3@gmail.com');
INSERT INTO users (email) VALUES ('phat4@gmail.com');
INSERT INTO users (email) VALUES ('phat5@gmail.com');
INSERT INTO users (email) VALUES ('phat6@gmail.com');
INSERT INTO users (email) VALUES ('phat7@gmail.com');
INSERT INTO users (email) VALUES ('phat8@gmail.com');
INSERT INTO users (email) VALUES ('phat9@gmail.com');

INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (1, 2, 1, 1);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (1, 3, 1, 1);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (1, 4, 1, 1);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (2, 3, 1, 2);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (2, 4, 1, 2);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (2, 5, 1, 2);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (3, 4, 1, 3);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (3, 5, 1, 3);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (3, 6, 1, 3);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (4, 5, 1, 4);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (4, 6, 1, 4);
INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
VALUES (4, 7, 1, 4);

INSERT INTO subcribers (requestor, target, status)
VALUES (1, 2, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (2, 1, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (3, 1, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (4, 1, 1);

INSERT INTO subcribers (requestor, target, status)
VALUES (2, 3, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (3, 2, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (4, 2, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (5, 2, 1);

INSERT INTO subcribers (requestor, target, status)
VALUES (3, 4, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (4, 3, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (5, 3, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (6, 3, 1);

INSERT INTO subcribers (requestor, target, status)
VALUES (4, 5, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (5, 4, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (6, 4, 1);
INSERT INTO subcribers (requestor, target, status)
VALUES (7, 4, 1);