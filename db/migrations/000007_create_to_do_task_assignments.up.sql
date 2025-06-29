CREATE TABLE to_do_task_assignments (
    to_do_list_detail_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (to_do_list_detail_id , user_id),
    FOREIGN KEY (to_do_list_detail_id) REFERENCES to_do_list_details(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);