create table public.restaurant_users
(
    id         uuid unique        not null,
    created_at timestampz         not null,
    updated_at timestampz,
    deleted_at timestampz,
    document   varchar(14) unique not null,
    email      varchar(50) unique not null,
    is_admin   boolean default false,

    constraint restaurant_users_pk
        PRIMARY KEY (id)
);

create unique index restaurant_users_document_index
    on public.restaurant_users using BTREE(document);