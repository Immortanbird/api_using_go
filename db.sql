CREATE TABLE IF NOT EXISTS `user` (
	`uid` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`user_name` VARCHAR(20) NOT NULL UNIQUE,
	`email` VARCHAR(128) NOT NULL UNIQUE,
	`pwd_hash` VARCHAR(255) NOT NULL,
	`gender` ENUM('M', 'F', 'X') NOT NULL,
	`birth` TIMESTAMP,
	`address` VARCHAR(255),
	`tel` VARCHAR(32),
	`created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);