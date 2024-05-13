CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SET TIMEZONE="Asia/Jakarta";

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL CHECK (LENGTH(name)>=5),
    password VARCHAR(255) NOT NULL,
    user_status INT NOT NULL,
    user_role VARCHAR(25) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL
);

CREATE TYPE cat_race AS ENUM('Persian','Maine Coon', 'Siamese', 'Ragdoll', 'Bengal', 'Sphynx', 'British Shorthair', 'Abyssinian','Scottish Fold', 'Birman');
CREATE TYPE cat_sex  AS ENUM('male', 'female');
CREATE TYPE match_status AS ENUM('approved', 'pending', 'rejected');

CREATE TABLE IF NOT EXISTS cats (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(30) NOT NULL CHECK (LENGTH(name)>=1 ),
    race cat_race NOT NULL,
    sex cat_sex NOT NULL,
    ageinmonth INT NOT NULL CHECK (ageInMonth >= 1 AND ageInMonth <= 120082),
    description VARCHAR(200) NOT NULL CHECK (LENGTH(description)>=1),   
    hasmatched BOOLEAN DEFAULT false NOT NULL ,
    imageurls TEXT[] NOT NULL CHECK (array_length(imageUrls, 1) >= 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS cat_matches (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    cat_issuer_id UUID NOT NULL,
    cat_match_id  UUID NOT NULL,
    message VARCHAR(120) NOT NULL CHECK (LENGTH(message)>=5),
    status match_status DEFAULT 'pending' NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (cat_issuer_id) REFERENCES cats (id),
    FOREIGN KEY (cat_match_id) REFERENCES cats (id)
);
-- Add indexes