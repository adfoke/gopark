-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    mail TEXT NOT NULL UNIQUE
);

-- Insert sample data
INSERT OR IGNORE INTO users (name, mail) VALUES
    ('Test User 1', 'test1@example.com'),
    ('Test User 2', 'test2@example.com');