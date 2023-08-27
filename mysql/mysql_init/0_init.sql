CREATE DATABASE IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS `user` (
    `private_user_id` INT AUTO_INCREMENT NOT NULL,
    `public_user_id` VARCHAR(36) NOT NULL,
    `user_name` VARCHAR(50) NOT NULL,
    `id` VARCHAR(50) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    `user_icon_id` VARCHAR(36),
    PRIMARY KEY (`private_user_id`)
);