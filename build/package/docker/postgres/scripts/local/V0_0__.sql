CREATE SCHEMA ads_campaing;

GRANT ALL PRIVILEGES ON DATABASE "ads-campaing-db" TO "postgres";

GRANT USAGE ON SCHEMA ads_campaing TO "postgres";
ALTER USER "postgres" SET search_path = 'ads_campaing';


SET SCHEMA 'ads_campaing';
ALTER DEFAULT PRIVILEGES
    IN SCHEMA ads_campaing
GRANT SELECT, UPDATE, INSERT, DELETE ON TABLES
    TO "postgres";

ALTER DEFAULT PRIVILEGES
    IN SCHEMA ads_campaing
GRANT USAGE ON SEQUENCES
    TO "postgres";