CREATE TABLE to_do_list (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    created_by INT NULL,
    updated_by INT NULL,
    deleted_by INT null,

    CONSTRAINT fk_todo_list_type FOREIGN KEY (type_id) REFERENCES to_do_types(id) ON DELETE CASCADE
);