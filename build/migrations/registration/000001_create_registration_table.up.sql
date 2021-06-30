create table if not exists registration
(
    id           serial                              not null
        constraint registration_pk
            primary key,
    email        varchar(256)                        not null,
    login        varchar(256)                        not null,
    code         int                                 not null,
    created_at   timestamp default current_timestamp not null,
    confirmed_at timestamp default null
);
create unique index if not exists registration_email_uindex on registration (email);
create unique index if not exists registration_login_uindex on registration (login);
