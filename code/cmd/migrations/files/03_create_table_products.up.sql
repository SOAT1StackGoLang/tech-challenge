create table public.lanchonete_products
(
    id          uuid           not null,
    created_at  timestamptz    not null,
    name        varchar(100)   not null,
    description varchar(200)   not null,
    category    varchar(30)    not null,
    price       numeric(10, 2) not null,

    constraint lanchonete_products_pk
        PRIMARY KEY (id)
);

create index lanchonete_products_category_index
    on public.lanchonete_products using hash (category);