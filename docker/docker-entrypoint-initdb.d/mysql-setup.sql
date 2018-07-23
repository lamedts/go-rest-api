CREATE DATABASE IF NOT EXISTS yay;
USE yay;

CREATE TABLE IF NOT EXISTS orders_origin (
    id INT PRIMARY KEY AUTO_INCREMENT,
    latitude FLOAT NOT NULL,
    longtitude FLOAT NOT NULL,
    order_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders_destination (
    id INT PRIMARY KEY AUTO_INCREMENT,
    latitude FLOAT NOT NULL,
    longtitude FLOAT NOT NULL,
    order_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    status ENUM('UNASSIGN', 'ASSIGNED') NOT NULL DEFAULT 'UNASSIGN',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    distance INT NOT NULL,
    origin_id INT NOT NULL,
    destination_id INT NOT NULL,
    FOREIGN KEY (origin_id)
        REFERENCES orders_origin (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (destination_id)
        REFERENCES orders_destination (id)
        ON UPDATE CASCADE ON DELETE CASCADE
);


