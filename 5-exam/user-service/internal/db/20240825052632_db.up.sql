create table users (
    user_id serial primary key,
    name varchar(64),
    email varchar(64),
    password text,
    age int
);