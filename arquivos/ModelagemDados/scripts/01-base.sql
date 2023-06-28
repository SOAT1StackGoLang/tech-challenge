DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'lanchonete') THEN
        CREATE DATABASE lanchonete;
    END IF;
END $$;

\c lanchonete