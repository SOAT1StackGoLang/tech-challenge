create table public.restaurant_users
(
    id uuid NOT NULL,
    created_at timestampz not null,
    updated_at timestampz,
    deleted_at timestampz,
    document varchar(14) not null,
    email varchar(50) not null,
    is_admin boolean default false,

    constraint restaurant_users_pk
        PRIMARY KEY (id)
);

create index restaurant_users_email_index
    on public.restaurant_users using BTREE(email);