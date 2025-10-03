-- PostgreSQL Database Setup Script
-- Run this as a PostgreSQL superuser (like postgres)

-- Create database
CREATE DATABASE currency_tracker;

-- Create user (optional - you can use existing postgres user)
CREATE USER currency_app WITH PASSWORD 'currency123';

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE currency_tracker TO currency_app;

-- Connect to the new database and grant schema permissions
\c currency_tracker;
GRANT ALL ON SCHEMA public TO currency_app;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO currency_app;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO currency_app;

-- The application will automatically create the currency_rates table
-- when it runs for the first time