-- Enable logging
ALTER SYSTEM SET log_destination = 'stderr';
ALTER SYSTEM SET log_statement = 'all';
ALTER SYSTEM SET log_min_duration_statement = 0;
ALTER SYSTEM SET log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h ';
ALTER SYSTEM SET log_rotation_age = '1d';
ALTER SYSTEM SET log_rotation_size = 0;

-- Enable statistics
ALTER SYSTEM SET track_activities = on;
ALTER SYSTEM SET track_counts = on;
ALTER SYSTEM SET track_io_timing = on;
ALTER SYSTEM SET track_functions = 'all';
ALTER SYSTEM SET track_activity_query_size = 2048;

-- Get configs
SELECT name, setting, source FROM pg_settings WHERE source != 'default';

-- Create users for msvc-orders, msvc-production, msvc-payments and password for each user for testing
CREATE USER msvc_orders WITH PASSWORD 'msvc_orders';
CREATE USER msvc_production WITH PASSWORD 'msvc_production';
CREATE USER msvc_payments WITH PASSWORD 'msvc_payments';

-- Create databases and give permission for msvc-orders(this will be the schema migrator must be admin), msvc-production, msvc-payments
CREATE DATABASE msvc_orders;
-- Grant admin permission to the database on msvc_orders and public schema
GRANT ALL PRIVILEGES ON DATABASE msvc_orders TO msvc_orders;
GRANT ALL PRIVILEGES ON SCHEMA public TO msvc_orders;



CREATE DATABASE msvc_production;
GRANT ALL PRIVILEGES ON DATABASE msvc_production TO msvc_production;
CREATE DATABASE msvc_payments;
GRANT ALL PRIVILEGES ON DATABASE msvc_payments TO msvc_payments;
