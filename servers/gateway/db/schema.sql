CREATE TABLE IF NOT EXISTS alerts (
    id INT not null auto_increment primary key,
    msg TEXT not null,
    deviceID INT not null,
    status BOOLEAN not null, -- sent or received
    created_at DATETIME not null DEFAULT GETDATE(),
    edited_at DATETIME -- can be null
);