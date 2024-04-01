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
    task_id    bigint      not null,
    user_id    bigint      not null,
    comment    text,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);

create table if not exists public.user
(
    id          bigserial primary key,
    group_id    bigint      not null,
    email       text        not null,
    password    text        not null,
    first_name  text        not null,
    last_name   text        not null,
    middle_name text,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz
);

create table if not exists public.group
(
    id         bigserial primary key,
    name       text        not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);

create table if not exists public.file
(
    id         bigserial primary key,
    name       text        not null,
    filename   text        not null,
    filepath   text        not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    task_id    bigint,
    answer_id  bigint
);