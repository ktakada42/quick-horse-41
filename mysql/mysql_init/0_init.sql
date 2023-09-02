CREATE DATABASE IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS `user` (
    `private_user_id` VARCHAR(26) NOT NULL,
    `public_user_id` VARCHAR(50) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    `user_name` VARCHAR(50) NOT NULL,
    `office_id` VARCHAR(26) NOT NULL,
    `user_icon_id` VARCHAR(26),
    `reg_date` DATETIME NOT NULL,
    `isAdmin` BOOLEAN NOT NULL,
    PRIMARY KEY (`private_user_id`)
);

CREATE TABLE IF NOT EXISTS `book` (
    `book_id` VARCHAR(26) NOT NULL,
    `isbn` VARCHAR(13),
    `title` VARCHAR(100) NOT NULL,
    `author` VARCHAR(100) NOT NULL,
    `publisher` VARCHAR(100),
    `publish_date` VARCHAR(8) NOT NULL,
    `cover_id` VARCHAR(26),
    PRIMARY KEY (`book_id`)
);

CREATE TABLE IF NOT EXISTS `office` (
    `office_id` VARCHAR(26) NOT NULL,
    `office_name` VARCHAR(50) NOT NULL,
    PRIMARY KEY (`office_id`)
);
