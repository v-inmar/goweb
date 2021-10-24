DROP DATABASE IF EXISTS `tododb`;
CREATE DATABASE IF NOT EXISTS `tododb`;
USE `tododb`;

--
-- Table structure for table `todo_model`
--

DROP TABLE IF EXISTS `todo_model`;
CREATE TABLE `todo_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `date_created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


--
-- Table structure for table `body_model`
--

DROP TABLE IF EXISTS `body_model`;
CREATE TABLE `body_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `value` text COLLATE utf8mb4_bin,
  `date_created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

--
-- Table structure for table `pid_model`
--

DROP TABLE IF EXISTS `pid_model`;
CREATE TABLE `pid_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `value` varchar(16) COLLATE utf8mb4_bin NOT NULL,
  `date_created` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `value_UNIQUE` (`value`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

--
-- Table structure for table `title_model`
--

DROP TABLE IF EXISTS `title_model`;
CREATE TABLE `title_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `value` varchar(128) COLLATE utf8mb4_bin NOT NULL,
  `date_created` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `value_UNIQUE` (`value`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

--
-- Table structure for table `todo_body_linker_model`
--

DROP TABLE IF EXISTS `todo_body_linker_model`;
CREATE TABLE `todo_body_linker_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `date_created` datetime NOT NULL,
  `date_updated` datetime NOT NULL,
  `todo_id` bigint NOT NULL,
  `body_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `body_id` (`body_id`),
  KEY `todo_id` (`todo_id`),
  CONSTRAINT `todo_body_linker_model_ibfk_1` FOREIGN KEY (`body_id`) REFERENCES `body_model` (`id`),
  CONSTRAINT `todo_body_linker_model_ibfk_2` FOREIGN KEY (`todo_id`) REFERENCES `todo_model` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



--
-- Table structure for table `todo_pid_linker_model`
--

DROP TABLE IF EXISTS `todo_pid_linker_model`;
CREATE TABLE `todo_pid_linker_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `date_created` datetime NOT NULL,
  `date_updated` datetime NOT NULL,
  `todo_id` bigint NOT NULL,
  `pid_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `pid_id` (`pid_id`),
  KEY `todo_id` (`todo_id`),
  CONSTRAINT `todo_pid_linker_model_ibfk_1` FOREIGN KEY (`pid_id`) REFERENCES `pid_model` (`id`),
  CONSTRAINT `todo_pid_linker_model_ibfk_2` FOREIGN KEY (`todo_id`) REFERENCES `todo_model` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Table structure for table `todo_title_linker_model`
--

DROP TABLE IF EXISTS `todo_title_linker_model`;
CREATE TABLE `todo_title_linker_model` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `date_created` datetime NOT NULL,
  `date_updated` datetime NOT NULL,
  `todo_id` bigint NOT NULL,
  `title_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `title_id` (`title_id`),
  KEY `todo_id` (`todo_id`),
  CONSTRAINT `todo_title_linker_model_ibfk_1` FOREIGN KEY (`title_id`) REFERENCES `title_model` (`id`),
  CONSTRAINT `todo_title_linker_model_ibfk_2` FOREIGN KEY (`todo_id`) REFERENCES `todo_model` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
