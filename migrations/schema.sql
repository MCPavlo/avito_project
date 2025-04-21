CREATE TABLE IF NOT EXISTS users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       role VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS pvz (
                     id SERIAL PRIMARY KEY,
                     city VARCHAR(255) NOT NULL,
                     registration TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS receptions (
                            id SERIAL PRIMARY KEY,
                            pvz_id INT NOT NULL,
                            created_at TIMESTAMP NOT NULL,
                            status VARCHAR(50) NOT NULL,
                            FOREIGN KEY (pvz_id) REFERENCES pvz(id)
);

CREATE TABLE IF NOT EXISTS goods (
                       id SERIAL PRIMARY KEY,
                       type VARCHAR(255) NOT NULL,
                       received_at TIMESTAMP NOT NULL,
                       reception_id INT NOT NULL,
                       FOREIGN KEY (reception_id) REFERENCES receptions(id)
);