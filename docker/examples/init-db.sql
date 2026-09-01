-- Evolution GO database initialization script

-- Create authentication database
CREATE DATABASE evogo_auth;

-- Create user data database
CREATE DATABASE evogo_users;

-- Confirmation message
SELECT 'Databases evogo_auth and evogo_users created successfully!' as message;

