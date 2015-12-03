CREATE TABLE users (
  name             varchar(255),
  email            varchar(36) PRIMARY KEY,
  password         varchar(60),
  activated        boolean,
  sam              varchar(20)
);
