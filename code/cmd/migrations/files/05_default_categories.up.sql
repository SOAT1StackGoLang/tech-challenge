insert into public.lanchonete_categories (id, created_at, name)
values ('9764bd96-3bcf-11ee-be56-0242ac120002', now(), 'Lanche');

insert into public.lanchonete_categories (id, created_at, name)
values ('a0424802-3bcf-11ee-be56-0242ac120002', now(), 'Acompanhamento');

insert into public.lanchonete_categories (id, created_at, name)
values ('a557b0c0-3bcf-11ee-be56-0242ac120002', now(), 'Bebida');

insert into public.lanchonete_categories (id, created_at, name)
values ('b1bef3fa-3bcf-11ee-be56-0242ac120002', now(), 'Sobremesa');

create unique index lanchonete_categories_name_index
    on public.lanchonete_categories using BTREE (name);