CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    Account VARCHAR(255) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    created_by INT NULL,
    updated_by INT NULL,
    deleted_by INT null
);