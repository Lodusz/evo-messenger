-- auth 
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- chat 
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_messages_id_asc ON messages(id ASC);

INSERT INTO users (username, password_hash) VALUES 
('Тимур', 'fake_hash_1'), 
('Юра', 'fake_hash_2')
ON CONFLICT (username) DO NOTHING;

INSERT INTO messages (username, content) VALUES
('Тимур', 'Всем привет! Тестим мессенджер Максика'),
('Юра', 'Ку! Всё работает чётко');