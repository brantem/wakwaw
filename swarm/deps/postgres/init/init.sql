CREATE TABLE users (
  id varchar(255) NOT NULL,
  name varchar(255) NOT NULL,
  PRIMARY KEY (id)
);

INSERT INTO users (id, name) VALUES ('user-1', 'User 1')