CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username) VALUES ('Тимур'), ('Юра'), ('Элина'), ('Камиль');
INSERT INTO messages (user_id, content) VALUES 
(1, 'Всем привет! Тестируем новый чат?'),
(2, 'Да, вроде работает.'),
(3, 'Ага, сообщения теперь сохраняются'),
(4, 'Отлично, щас тоже зайду с ноута.');