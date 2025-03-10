-- статус: pending, sent, executed, error, online, offline
CREATE TABLE IF NOT EXISTS statuses (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT,
    code TEXT NOT NULL UNIQUE,
    context TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO statuses ("name", code, context) VALUES('offline', 'OFFLINE', 'device');
INSERT INTO statuses ("name", code, context) VALUES('online', 'ONLINE', 'device');
INSERT INTO statuses ("name", code, context) VALUES('pending', 'PENDING', 'commands');
INSERT INTO statuses ("name", code, context) VALUES('sent', 'SENT', 'commands');
INSERT INTO statuses ("name", code, context) VALUES('executed', 'EXECUTED', 'commands');
INSERT INTO statuses ("name", code, context) VALUES('error', 'ERROR', 'commands');

CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT,
    code TEXT NOT NULL UNIQUE,
    description TEXT
);

INSERT INTO roles ("name", code) VALUES('observer', 'OBSERVER');
INSERT INTO roles ("name", code) VALUES('observerplus', 'OBSERVER_PLUS');
INSERT INTO roles ("name", code) VALUES('admin', 'ADMIN');

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role_code TEXT REFERENCES roles(code),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO users (id, username, email, password_hash, role_id) VALUES('a7c4265d-545c-404f-a1ef-daf4af2dfb12', 'admin', 'admin@example.com', '$2y$10$yUilfuMXj6ZIldGxdQpJUOHey8pvinWGBHWXGs8zcZGENMw1J6z2C', NULL);

-- Таблица устройств (устройства, которыми управляет система)
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_identifier TEXT NOT NULL UNIQUE,
    group_id TEXT REFERENCES device_groups(id) ON DELETE SET NULL,
    description TEXT,
    status TEXT REFERENCES statuses(code),
    last_seen TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица метрик, получаемых от клиентов
CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_id TEXT REFERENCES devices(id) on delete cascade,
    public_ip TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    hostname TEXT,
    os_info TEXT,
    disk_total BIGINT,
    disk_used BIGINT,
    disk_free BIGINT,
    memory_total BIGINT,
    memory_used BIGINT,
    memory_available BIGINT,
    process_count INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    -- В будущем можно реализовать партиционирование по created_at для таймсерийных данных
);

-- Таблица команд, отправленных устройствам
CREATE TABLE IF NOT EXISTS commands (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_id TEXT REFERENCES devices(id) on delete cascade,
    command_type TEXT NOT NULL,
    user_id TEXT REFERENCES users(id),
    status TEXT REFERENCES statuses(code) ON DELETE SET NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    executed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS device_groups (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS applications (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    version TEXT,
    app_type TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS device_applications (
    device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    installed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (device_id, application_id)
);