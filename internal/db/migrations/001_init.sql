-- Create the berries database if it doesn't exist
CREATE DATABASE IF NOT EXISTS poke_app;

-- Switch to the berries database
USE poke_app;

-- Create the berries table
CREATE TABLE IF NOT EXISTS `berries` (
                                      id INT AUTO_INCREMENT PRIMARY KEY,
                                      name VARCHAR(255) NOT NULL,
                                      url VARCHAR(255) NOT NULL,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);