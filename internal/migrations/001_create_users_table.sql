-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mail VARCHAR(255) NOT NULL UNIQUE
);

-- Insert sample data
INSERT INTO users (name, mail) VALUES
    ('Test User 1', 'test1@example.com'),
    ('Test User 2', 'test2@example.com')
ON CONFLICT (mail) DO NOTHING;
