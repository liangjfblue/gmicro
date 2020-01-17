create schema if not exists db_gmicro collate utf8mb4;

create table if not exists tb_article
(
	id int unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	uid varchar(100) null,
	content text null,
	article_info_id int unsigned not null
);

create index idx_tb_article_deleted_at
	on tb_article (deleted_at);

create table if not exists tb_article_info
(
	id int unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	title varchar(100) null,
	topic varchar(100) null,
	author varchar(100) null,
	is_original tinyint null
);

create index idx_tb_article_info_deleted_at
	on tb_article_info (deleted_at);

create table if not exists tb_comment
(
	id int unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	uid varchar(100) null,
	article_id int unsigned not null,
	comment varchar(200) null,
	from_id int unsigned not null,
	to_id int unsigned not null
);

create index idx_tb_comment_deleted_at
	on tb_comment (deleted_at);

create table if not exists tb_money
(
	id int unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	uid varchar(100) null,
	coin int null,
	constraint uix_tb_money_uid
		unique (uid)
);

create index idx_tb_money_deleted_at
	on tb_money (deleted_at);

create table if not exists tb_user
(
	id int unsigned auto_increment
		primary key,
	created_at datetime null,
	updated_at datetime null,
	deleted_at datetime null,
	uid varchar(100) null,
	username varchar(100) null,
	password varchar(80) null,
	age int null,
	address longtext null,
	is_available tinyint null,
	last_login datetime null,
	login_ip varchar(20) null,
	mid int unsigned null,
	constraint uix_tb_user_uid
		unique (uid),
	constraint uix_tb_user_username
		unique (username)
);

create index idx_tb_user_deleted_at
	on tb_user (deleted_at);

