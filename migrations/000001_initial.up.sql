create table if not exists public.task
(
    id             bigserial primary key,
    group_id       bigint      not null,
    title          text        not null,
    text           text,
    cost           bigint      not null,
    effective_from timestamptz,
    effective_till timestamptz,
    created_at     timestamptz not null default now(),
    updated_at     timestamptz
);

create table if not exists public.answer
(
    id         bigserial primary key,
    group_id   bigint      not null,
    comment    text,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);