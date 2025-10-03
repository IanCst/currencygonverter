# PostgreSQL Database Setup Guide

## Quick Setup Steps

### 1. Install PostgreSQL (if not already installed)
- Download from: https://www.postgresql.org/download/windows/
- During installation, remember the password you set for the `postgres` user

### 2. Create Database and User
Run these commands in PostgreSQL command line (psql) or pgAdmin:

```sql
-- Connect as postgres user first
-- Then run these commands:

CREATE DATABASE currency_tracker;
CREATE USER currency_app WITH PASSWORD 'currency123';
GRANT ALL PRIVILEGES ON DATABASE currency_tracker TO currency_app;

-- Connect to the new database
\c currency_tracker;
GRANT ALL ON SCHEMA public TO currency_app;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO currency_app;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO currency_app;
```

### 3. Test Connection
You can test the connection with:
```bash
psql -h localhost -U currency_app -d currency_tracker
# Enter password: currency123
```

## Default Configuration

The app uses these default settings (from .env file):

- **Host**: localhost
- **Port**: 5432
- **Database**: currency_tracker
- **User**: currency_app
- **Password**: currency123

## Alternative: Use Existing PostgreSQL User

If you want to use your existing `postgres` user instead:

1. Edit the `.env` file:
   ```
   DB_USER=postgres
   DB_PASSWORD=your_postgres_password
   ```

2. Create only the database:
   ```sql
   CREATE DATABASE currency_tracker;
   ```

## Troubleshooting

1. **Authentication failed**: Make sure the user and password are correct
2. **Connection refused**: Check if PostgreSQL service is running
3. **Database doesn't exist**: Run the CREATE DATABASE command first
4. **Permission denied**: Make sure the user has proper privileges

The application will work without a database connection and will simply show warnings in the logs.