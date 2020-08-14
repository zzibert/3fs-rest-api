
-- CREATE TABLE groups (
--   id serial PRIMARY KEY,
--   name varchar(50) UNIQUE NOT NULL
-- );

-- CREATE TABLE users (
--   id serial PRIMARY KEY,
--   name varchar(50) NOT NULL,
--   password varchar(50) NOT NULL,
--   email varchar(255) UNIQUE NOT NULL,
--   group_id integer NOT NULL references groups(id)
-- );