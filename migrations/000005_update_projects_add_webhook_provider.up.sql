ALTER TABLE projects
ADD COLUMN webhook_provider VARCHAR(255) DEFAULT 'generic';