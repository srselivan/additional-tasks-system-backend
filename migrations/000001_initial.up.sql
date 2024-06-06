create table if not exists public.task
(
    id             bigserial primary key,
    title          text        not null,
    text           text,
    created_by     bigint      not null references public."user" (id),
    effective_from timestamptz,
    effective_till timestamptz,
    created_at     timestamptz not null default now(),
    updated_at     timestamptz
);

create table if not exists public.task_links
(
    id       bigserial primary key,
    user_id  bigint references public."user" (id),
    group_id bigint references public."group" (id),
    task_id  bigint not null references public.task (id),

    constraint unique_user_task unique (user_id, task_id),
    constraint unique_group_task unique (group_id, task_id),
    constraint not_null_user_group check ((user_id is not null or group_id is not null))
);

create table if not exists public.answer
(
    id         bigserial primary key,
    task_id    bigint      not null references public.task (id),
    user_id    bigint      not null references public."user" (id),
    comment    text,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);

create table if not exists public.user
(
    id          bigserial primary key,
    group_id    bigint references public.group (id),
    role_id     bigint      not null references public.roles (id),
    email       text unique not null,
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
    task_id    bigint references public.task (id),
    answer_id  bigint references public.answer (id)
);

create table if not exists public.user_token
(
    user_id bigint not null references public."user" (id) unique,
    refresh text   not null
);

create table if not exists public.roles
(
    id   bigserial primary key,
    name text not null
);

insert into public.roles (name)
values ('administrator'),
       ('teacher'),
       ('student');

create table if not exists public.mark
(
    id         bigserial primary key,
    mark       integer     not null,
    comment    text,
    answer_id  bigint      not null references public.answer (id),
    created_at timestamptz not null default now(),
    updated_at timestamptz
);