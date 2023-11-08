CREATE DATABASE IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS `office` (
    `office_id` VARCHAR(26) NOT NULL UNIQUE,
    `office_name` VARCHAR(50) NOT NULL,
    PRIMARY KEY (`office_id`)
);

CREATE TABLE IF NOT EXISTS `image` (
    `image_id` VARCHAR(26) NOT NULL UNIQUE,
    `image_name` VARCHAR(120) NOT NULL,
    `path` VARCHAR(1024) NOT NULL,
    PRIMARY KEY (`image_id`)
);

CREATE TABLE IF NOT EXISTS `user` (
    `user_id` VARCHAR(26) NOT NULL UNIQUE,
    `password` VARCHAR(100) NOT NULL,
    `user_name` VARCHAR(50) NOT NULL,
    `office_id` VARCHAR(26) NOT NULL,
    `user_icon_id` VARCHAR(26) NOT NULL DEFAULT '',
    `reg_date` DATETIME NOT NULL,
    `isAdmin` BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (`user_id`),
    FOREIGN KEY (office_id) REFERENCES office(office_id),
    FOREIGN KEY (user_icon_id) REFERENCES image(image_id)
);

CREATE TABLE IF NOT EXISTS `book_master` (
    `book_id` VARCHAR(26) NOT NULL UNIQUE,
    `isbn` VARCHAR(13) UNIQUE,
    `title` VARCHAR(100) NOT NULL,
    `author` VARCHAR(100) NOT NULL,
    `publisher` VARCHAR(100) NOT NULL DEFAULT '',
    `publish_date` VARCHAR(8) NOT NULL,
    `cover_id` VARCHAR(26) NOT NULL DEFAULT '',
    PRIMARY KEY (`book_id`),
    FOREIGN KEY (cover_id) REFERENCES image(image_id)
);

CREATE TABLE IF NOT EXISTS `office_book` (
    `office_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    `office_book_id` INT NOT NULL,
    `reg_date` DATETIME NOT NULL,
    `available` BOOLEAN NOT NULL DEFAULT TRUE,
    `metadata` VARCHAR(50) NOT NULL DEFAULT '',
    PRIMARY KEY (`office_id`, `book_id`, `office_book_id`),
    INDEX `idx_office_book_id` (`office_book_id`),
    FOREIGN KEY (office_id) REFERENCES office(office_id),
    FOREIGN KEY (book_id) REFERENCES book_master(book_id)
);

CREATE TABLE IF NOT EXISTS `borrow_history` (
    `book_id` VARCHAR(26) NOT NULL,
    `history_id` INT NOT NULL,
    `user_id` VARCHAR(26) NOT NULL,
    `borrow_date` DATETIME NOT NULL,
    `return_date` DATETIME NOT NULL,
    PRIMARY KEY (`book_id`, `history_id`),
    FOREIGN KEY (book_id) REFERENCES book_master(book_id),
    FOREIGN KEY (user_id) REFERENCES user(user_id)
);

CREATE TABLE IF NOT EXISTS `borrowed_book` (
    `office_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    `office_book_id` INT NOT NULL,
    `user_id` VARCHAR(26) NOT NULL,
    `borrow_date` DATETIME NOT NULL,
    PRIMARY KEY (`office_id`, `book_id`, `office_book_id`),
    FOREIGN KEY (office_id) REFERENCES office(office_id),
    FOREIGN KEY (book_id) REFERENCES book_master(book_id),
    FOREIGN KEY (office_book_id) REFERENCES office_book(office_book_id),
    FOREIGN KEY (user_id) REFERENCES user(user_id)
);

CREATE TABLE IF NOT EXISTS `review` (
    `book_id` VARCHAR(26) NOT NULL,
    `review_id` INT NOT NULL,
    `user_id` VARCHAR(26) NOT NULL,
    `rating` TINYINT NOT NULL DEFAULT -1,
    `review` VARCHAR(500) NOT NULL DEFAULT '',
    `reg_date` DATETIME NOT NULL,
    PRIMARY KEY (`book_id`, `user_id`),
    FOREIGN KEY (book_id) REFERENCES book_master(book_id),
    FOREIGN KEY (user_id) REFERENCES user(user_id)
);