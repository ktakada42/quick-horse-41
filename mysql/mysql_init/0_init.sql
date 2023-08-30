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
