CREATE TABLE IF NOT EXISTS users (
   id serial PRIMARY KEY,
   email TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS friends (
   user_one_id INTEGER NOT NULL,
   user_two_id INTEGER NOT NULL,
   status INTEGER NOT NULL,
   user_action_id INTEGER NOT NULL,
   PRIMARY KEY (user_one_id, user_two_id)
);

CREATE TABLE IF NOT EXISTS subcribers (
   id SERIAL PRIMARY KEY NOT NULL,
   requestor INT NOT NULL,
   target INT NOT NULL,
   status INTEGER NOT NULL
);