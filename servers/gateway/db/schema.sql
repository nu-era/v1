CREATE TABLE IF NOT EXISTS alerts (
    id INT not null auto_increment primary key,
    msg TEXT not null,
    device_name TEXT not null,
    status BOOLEAN not null, -- sent or received
    created_at DATETIME not null DEFAULT GETDATE(),
    edited_at DATETIME, -- can be null
    send_time DATETIME not null, 
    receive_time DATETIME not null 
);