create table public.lanchonete_orders
(
    id         uuid           not null,
    created_at timestamptz    not null,
    updated_at timestamptz,
    deleted_at timestamptz,
    user_id    uuid           not null,
    products   json           not null,
    price      numeric(10, 2) not null,
    status     int default 0  not null,
    payment_id uuid,

    constraint lanchonete_orders_pk
        PRIMARY KEY (id)
);

alter table public.lanchonete_orders
    add constraint fk_order_user_id
        foreign key (user_id)
            references public.lanchonete_users (id);