CREATE TABLE users (
    user_id serial primary key not null unique ,
    username varchar(255) not null unique ,
    email varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE auth (
    accessCode varchar(20) not null ,
    refreshToken varchar(255) not null
);