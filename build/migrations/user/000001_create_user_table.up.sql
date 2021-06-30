create table if not exists "user"
(
    id         serial                                 not null
        constraint user_pk primary key,
    email      varchar(256)                           not null,
    login      varchar(256)                           not null,
    bio        varchar(160) default ''                not null,
    created_at timestamp    default current_timestamp not null
);
create unique index if not exists user_email_uindex on "user" (email);
create unique index if not exists user_login_uindex on "user" (login);
