create table if not exists products (
  id bigserial primary key,
  name varchar(255) not null,
  "type" varchar(10) not null,
  price integer not null,
  stocks integer not null,
  sold integer default 0 not null,
  created_time timestamp with time zone default now() not null,
  updated_time timestamp with time zone default now() not null
);
