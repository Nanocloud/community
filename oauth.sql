CREATE TABLE oauth_clients (
  id      serial PRIMARY KEY,
  name    varchar(255) UNIQUE,
  key     varchar(255) UNIQUE,
  secret  varchar(255)
);

CREATE TABLE oauth_access_tokens (
  id                serial PRIMARY KEY,
  token             varchar(255) UNIQUE,
  oauth_client_id   integer REFERENCES oauth_clients (id),
  user_id           varchar(255)
);
