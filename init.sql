-- Auth Service
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--Chat Service 
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- INDEX
CREATE INDEX IF NOT EXISTS idx_messages_id_asc ON messages(id ASC);

-- Тестовые данные для стартовой истории
INSERT INTO users (username) VALUES ('Тимур'), ('Юра')
ON CONFLICT (username) DO NOTHING;

INSERT INTO messages (username, content) VALUES 
('Тимур', 'Всем привет! Тестим мессенджер Максика'),
('Юра', 'Ку! Всё работает чётко');