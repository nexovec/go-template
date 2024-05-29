CREATE SCHEMA IF NOT EXISTS rbac;
-- users can have both roles and permissions, roles have permissions, permissions are something atomic

-- Users
CREATE TABLE IF NOT EXISTS rbac.users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    first_name VARCHAR,
    last_name VARCHAR,
    username VARCHAR NOT NULL,
    password VARCHAR,
    email VARCHAR NOT NULL,
    description VARCHAR,
    account_locked BOOLEAN NOT NULL DEFAULT FALSE,
    path_to_avatar VARCHAR,

    permissions VARCHAR[]
);

-- UserLogins
-- NOTE: multiple simultaneous logins ARE allowed
CREATE TABLE IF NOT EXISTS rbac.user_logins (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    terminated_at TIMESTAMP DEFAULT NULL,
    user_id INT NOT NULL,
    successful BOOLEAN DEFAULT FALSE NOT NULL,
    ip_address VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,

    FOREIGN KEY (user_id) REFERENCES rbac.users (id)
);

-- Roles
CREATE TABLE IF NOT EXISTS rbac.roles (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    name VARCHAR NOT NULL,
    description VARCHAR,

    permissions VARCHAR[]
);

-- PvtUserRole
CREATE TABLE IF NOT EXISTS rbac.pvt_user_role (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES rbac.users (id),
    FOREIGN KEY (role_id) REFERENCES rbac.roles (id)
);


-- Procedures
CREATE OR REPLACE FUNCTION rbac.insert_user(
    p_username VARCHAR DEFAULT NULL,
    p_email VARCHAR DEFAULT NULL,
    p_password VARCHAR DEFAULT NULL
)
RETURNS INT AS $$
DECLARE
    user_id INT;
BEGIN
    IF EXISTS (
        SELECT 1
        FROM users
        WHERE (username = p_username OR email = p_email) AND deleted_at IS NULL
    ) THEN
        RAISE EXCEPTION 'User already exists';
    ELSE
        INSERT INTO users (username, email, password)
        VALUES (p_username, p_email, p_password)
        RETURNING id INTO user_id;
    END IF;
    RETURN user_id;
END;
$$
LANGUAGE plpgsql;