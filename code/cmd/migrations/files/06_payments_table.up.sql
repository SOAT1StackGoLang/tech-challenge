create table public.lanchonete_payments
(
    id         uuid           not null,
    created_at timestamptz    not null,
    updated_at timestamptz,
    order_id   uuid unique    not null,
    value      numeric(10, 2) not null,
    status     varchar(20)    not null,

    constraint lanchonete_payments_pk
        PRIMARY KEY (id)
);

alter table public.lanchonete_payments
    add constraint fk_payment_order_id
        foreign key (order_id)
            references public.lanchonete_orders (id);