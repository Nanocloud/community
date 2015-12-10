CREATE TABLE users (
  id               varchar(36) PRIMARY KEY
  first_name       varchar(36),
  last_name        varchar(36),
  email            varchar(36) UNIQUE,
  password         varchar(60),
  is_admin         boolean,
  activated        boolean
);
