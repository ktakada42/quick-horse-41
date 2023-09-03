CREATE DATABASE IF NOT EXISTS app;

CREATE TABLE IF NOT EXISTS `image` (
    `image_id` VARCHAR(26) NOT NULL,
    `path` VARCHAR(100) NOT NULL,
    PRIMARY KEY (`image_id`)
);

CREATE TABLE IF NOT EXISTS `user` (
    `private_user_id` VARCHAR(26) NOT NULL,
    `public_user_id` VARCHAR(50) NOT NULL UNIQUE,
    `password` VARCHAR(100) NOT NULL,
    `user_name` VARCHAR(50) NOT NULL,
    `office_id` VARCHAR(26) NOT NULL,
    `user_icon_id` VARCHAR(26),
    `reg_date` DATETIME NOT NULL,
    `isAdmin` BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (`private_user_id`),
    FOREIGN KEY (user_icon_id) REFERENCES image(image_id)
);

CREATE TABLE IF NOT EXISTS `book` (
    `book_id` VARCHAR(26) NOT NULL,
    `isbn` VARCHAR(13) UNIQUE,
    `title` VARCHAR(100) NOT NULL,
    `author` VARCHAR(100) NOT NULL,
    `publisher` VARCHAR(100),
    `publish_date` VARCHAR(8) NOT NULL,
    `cover_id` VARCHAR(26),
    PRIMARY KEY (`book_id`),
    FOREIGN KEY (cover_id) REFERENCES image(image_id)
);

CREATE TABLE IF NOT EXISTS `office` (
    `office_id` VARCHAR(26) NOT NULL,
    `office_name` VARCHAR(50) NOT NULL,
    PRIMARY KEY (`office_id`)
);

CREATE TABLE IF NOT EXISTS `office_book` (
    `office_book_id` VARCHAR(26) NOT NULL,
    `office_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    `reg_date` DATETIME NOT NULL,
    `available` BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (`office_book_id`),
    FOREIGN KEY (office_id) REFERENCES office(office_id),
    FOREIGN KEY (book_id) REFERENCES book(book_id)
);

CREATE TABLE IF NOT EXISTS `borrow_history` (
    `history_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    `public_user_id` VARCHAR(50) NOT NULL,
    `borrow_date` DATETIME NOT NULL,
    `return_date` DATETIME NOT NULL,
    PRIMARY KEY (`history_id`),
    FOREIGN KEY (book_id) REFERENCES book(book_id),
    FOREIGN KEY (public_user_id) REFERENCES user(public_user_id)
);

CREATE TABLE IF NOT EXISTS `borrowed_book` (
    `office_book_id` VARCHAR(26) NOT NULL,
    `public_user_id` VARCHAR(50) NOT NULL,
    `borrow_date` DATETIME NOT NULL,
    PRIMARY KEY (`office_book_id`),
    FOREIGN KEY (office_book_id) REFERENCES office_book(office_book_id),
    FOREIGN KEY (public_user_id) REFERENCES user(public_user_id)
);

CREATE TABLE IF NOT EXISTS `request` (
    `request_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    `public_user_id` VARCHAR(50) NOT NULL,
    `comment` VARCHAR(500),
    `reg_date` DATETIME NOT NULL,
    PRIMARY KEY (`request_id`),
    FOREIGN KEY (book_id) REFERENCES book(book_id),
    FOREIGN KEY (public_user_id) REFERENCES user(public_user_id)
);

CREATE TABLE IF NOT EXISTS `requested_book` (
    `book_id` VARCHAR(26) NOT NULL,
    `processed` BOOLEAN NOT NULL,
    PRIMARY KEY (`book_id`),
    FOREIGN KEY (book_id) REFERENCES book(book_id)
);

CREATE TABLE IF NOT EXISTS `review` (
    `review_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    `public_user_id` VARCHAR(50) NOT NULL,
    `rating` TINYINT,
    `review` VARCHAR(500),
    `reg_date` DATETIME NOT NULL,
    PRIMARY KEY (`review_id`),
    FOREIGN KEY (book_id) REFERENCES book(book_id),
    FOREIGN KEY (public_user_id) REFERENCES user(public_user_id)
);

CREATE TABLE IF NOT EXISTS `tag` (
    `tag_id` VARCHAR(26) NOT NULL,
    `description` VARCHAR(20),
    PRIMARY KEY (`tag_id`)
);

CREATE TABLE IF NOT EXISTS `book_tag` (
    `book_tag_id` VARCHAR(26) NOT NULL,
    `tag_id` VARCHAR(26) NOT NULL,
    `book_id` VARCHAR(26) NOT NULL,
    PRIMARY KEY (`book_tag_id`),
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id),
    FOREIGN KEY (book_id) REFERENCES book(book_id)
);
