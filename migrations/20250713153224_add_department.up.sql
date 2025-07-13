CREATE TYPE departments AS ENUM (
    'Engineering',
    'Marketing',
    'Sales',
    'Human Resources',
    'Finance',
    'Operations',
    'Customer Support',
    'Product Management',
    'Quality Assurance',
    'IT/Infrastructure',
    'Legal',
    'Other'
);

ALTER TABLE tickets
ADD COLUMN department departments NOT NULL;