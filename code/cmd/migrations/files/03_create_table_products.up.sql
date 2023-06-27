create table public.lanchonete_categories
(
    id         uuid        not null,
    created_at timestamptz not null,
    updated_at timestamptz,
    name       varchar(40) not null,

    constraint lanchonete_categories_pk
        PRIMARY KEY (id)
);

create table public.lanchonete_products
(
    id          uuid           not null,
    created_at  timestamptz    not null,
    updated_at  timestamptz,
    category_id uuid           not null,
    name        varchar(100)   not null,
    description varchar(200)   not null,
    price       numeric(10, 2) not null,

    constraint lanchonete_products_pk
        PRIMARY KEY (id)
);

ALTER TABLE public.lanchonete_products
    ADD CONSTRAINT fk_product_category_id
        FOREIGN KEY (category_id)
            REFERENCES public.lanchonete_categories (id);