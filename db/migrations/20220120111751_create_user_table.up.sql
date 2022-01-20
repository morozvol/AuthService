create table users
(
    id                 serial
        constraint users_pk
            primary key,
    email              text not null,
    encrypted_password text not null,
    sold               text not null
);

alter table users
    owner to postgres;

create unique index users_email_uindex
    on users (email);

