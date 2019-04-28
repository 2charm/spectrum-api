create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(128) not null unique,
    pass_hash varchar(256) not null,
    user_name varchar(256) not null unique,
    first_name varchar(64) not null,
    last_name varchar(128) not null,
    photo_url varchar(128) not null
);

create table if not exists sign_in (
    user_id int not null,
    attempt_time datetime not null,
    client_ip varchar(128) not null
);
