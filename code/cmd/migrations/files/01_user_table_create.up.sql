create table public.lanchonete_users
(
    id         uuid               not null,
    created_at timestamptz        not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    name       varchar(100)       not null,
    document   varchar(14) unique not null,
    email      varchar(50) unique not null,
    is_admin   boolean default false,

    constraint lanchonete_users_pk
        PRIMARY KEY (id)
);

create unique index lanchonete_users_document_index
    on public.lanchonete_users using BTREE (document);