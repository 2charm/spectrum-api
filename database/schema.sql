create table if not exists users (
    user_id int not null auto_increment primary key,
    email varchar(128) not null unique,
    pass_hash varchar(256) not null,
    user_name varchar(256) not null unique,
    first_name varchar(64) not null,
    last_name varchar(128) not null
);

create table if not exists sign_in (
    user_id int not null,
    attempt_time datetime not null,
    client_ip varchar(128) not null
);

create table if not exists categories (
    category_id int not null auto_increment primary key,
    category_name varchar(128) not null unique
);

insert into categories
values(1, 'sports');
insert into categories
values(2, 'health');
insert into categories
values(3, 'business');
insert into categories
values(4, 'entertainment');
insert into categories
values(5, 'science');
insert into categories
values(6, 'technology');
insert into categories
values(7, 'headline');

create table if not exists articles (
    user_id int not null,
    category_id int not null,
    read_on datetime not null
);

create table if not exists sources (
    source_id int not null auto_increment primary key,
    source_name varchar(128) not null unique
);
