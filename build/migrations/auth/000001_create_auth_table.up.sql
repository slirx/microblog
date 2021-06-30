create table if not exists auth
(
    id         serial                              not null
        constraint auth_pk
            primary key,
    user_id    int                                 not null,
    email      varchar(256)                        not null,
    login      varchar(256)                        not null,
    password   varchar(60)                         not null,
    created_at timestamp default current_timestamp not null
);
create unique index if not exists auth_email_uindex on auth (email);
create unique index if not exists auth_login_uindex on auth (login);
