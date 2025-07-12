CREATE TYPE categories AS ENUM ('software', 'hardware', 'login', 'other');
CREATE TYPE priorities AS ENUM ('low', 'medium', 'high', 'critical');

CREATE TABLE IF NOT EXISTS tickets (
        id BIGSERIAL PRIMARY KEY,
        full_name VARCHAR(255),
        email_address VARCHAR(255),
        issue_category categories NOT NULL,
        priority priorities NOT NULL,
        issue_summary TEXT,
        detailed_description TEXT,
        created_at TIMESTAMPTZ NOT NULL,
        updated_at TIMESTAMPTZ NOT NULL
    );