create database cmdb ;
use cmdb;
create table engine_room (
    id int primary key AUTO_INCREMENT,
    city varchar(255) NOT NULL,
    engine_room varchar(30) NOT NULL ,
    cabinet varchar(30) NOT NULL,
    create_time DATETIME ,
    update_time DATETIME ,
    cabinet_number varchar(30) NOT NULL
    )  ENGINE = InnoDB AUTO_INCREMENT = 574 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = DYNAMIC ;

create table server (
    id int  primary key AUTO_INCREMENT,
    cabinet_number varchar(30) NOT NULL,
    name varchar(30) NOT NULL,
    model varchar(30) NOT NULL,
    label varchar(30) NOT NULL,
    ipaddresss varchar(30) NOT NULL,
    cpu    varchar(30) NOT NULL,
    memory varchar(30) NOT NULL,
    disk   varchar(30) NOT NULL,
    create_time DATETIME ,
    update_time DATETIME ,
    location    varchar(30) NOT NULL
)  ENGINE = InnoDB AUTO_INCREMENT = 574 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = DYNAMIC   ;

create table switch  (
       id int primary key AUTO_INCREMENT,
       name varchar(30) NOT NULL,
       model varchar(30) NOT NULL,
       location varchar(30) NOT NULL,
       ipaddresss varchar(30) NOT NULL,
       create_time DATETIME ,
       update_time DATETIME ,
       cabinet_number varchar(30) NOT NULL
)  ENGINE = InnoDB AUTO_INCREMENT = 574 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = DYNAMIC   ;

create table router (
    id int primary key AUTO_INCREMENT,
    name varchar(30) NOT NULL,
    model varchar(30) NOT NULL,
    location varchar(30) NOT NULL,
    ipaddresss varchar(30) NOT NULL,
    create_time DATETIME ,
    update_time DATETIME ,
    cabinet_number varchar(30) NOT NULL
)  ENGINE = InnoDB AUTO_INCREMENT = 574 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = DYNAMIC ;