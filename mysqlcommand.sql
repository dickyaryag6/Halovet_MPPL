CREATE TABLE `account` (
  `name` varchar(30) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(32) NOT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role` int(11) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) 

CREATE TABLE `appointment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `appointment_time` timestamp NOT NULL,
  `doctor_name` varchar(45) NOT NULL,
  `pet_owner_name` varchar(45) NOT NULL,
  `pet_type` varchar(45) NOT NULL,
  `complaint_description` varchar(45) DEFAULT NULL,
  `is_paid` tinyint(4) DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) 

CREATE TABLE `forum_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_title` varchar(45) NOT NULL,
  `category_desc` varchar(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) 

CREATE TABLE `forum_reply` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) DEFAULT NULL,
  `topic_id` int(11) DEFAULT NULL,
  `author` varchar(45) DEFAULT NULL,
  `content` varchar(45) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) 

CREATE TABLE `forum_topic` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) DEFAULT NULL,
  `title` varchar(45) DEFAULT NULL,
  `author` varchar(45) DEFAULT NULL,
  `content` varchar(45) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
)