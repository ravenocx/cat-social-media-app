-- Delete tables
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS cats CASCADE;
DROP TABLE IF EXISTS cat_matches;

DROP TYPE IF EXISTS cat_race;
DROP TYPE IF EXISTS cat_sex;
DROP TYPE IF EXISTS match_status;


DROP EXTENSION IF EXISTS "uuid-ossp";

