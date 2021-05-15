CREATE TABLE users (
  id bigserial not null primary key,
  email varchar not null unique,
  encrypted_password varchar not null
);

CREATE TABLE proxies (
  id serial not null primary key,
  host varchar not null ,
  login varchar not null unique,
  password varchar,
  isbad bool,
  isbusy bool,
  dataadd varchar,
  datalastuse varchar
);

