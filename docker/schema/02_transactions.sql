create table if not exists transactions (
  id bigserial primary key,
  user_id bigint not null,
  product_id bigint not null,
  item_amount integer not null,
  created_time timestamp with time zone default now() not null,
  updated_time timestamp with time zone default now() not null,
  constraint transactions_user_id_fk foreign key (user_id)
    references users(id),
  constraint transactions_product_id_fk foreign key (product_id)
    references products(id)
);
