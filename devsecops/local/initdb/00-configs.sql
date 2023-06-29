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