-- migrations/001_initial.sql

CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- статус: pending, sent, executed, error, online, offline
CREATE TABLE IF NOT EXISTS command_statuses (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT,
    code TEXT NOT NULL,
    context TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role_id TEXT REFERENCES roles(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица устройств (устройства, которыми управляет система)
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_identifier TEXT NOT NULL UNIQUE,
    user_id TEXT REFERENCES users(id),
    description TEXT,
    status TEXT,
    last_seen TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица метрик, получаемых от клиентов
CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    device_id TEXT REFERENCES devices(id),
    public_ip TEXT,
    hostname TEXT,
    os_info TEXT,
    disk_total BIGINT,
    disk_used BIGINT,
    disk_free BIGINT,
    memory_total BIGINT,
    memory_used BIGINT,
    memory_available BIGINT,
    process_count INTEGER,
    cpu_percent REAL,
    bytes_sent BIGINT,
    bytes_recv BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    -- В будущем можно реализовать партиционирование по created_at для таймсерийных данных
);

-- Таблица команд, отправленных устройствам
CREATE TABLE IF NOT EXISTS commands (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name ,
    code ,
    command_type TEXT NOT NULL,
    parameters JSONB,
    initiator INTEGER REFERENCES users(id),
    status TEXT REFERENCES command_statuses(code) NOT NULL DEFAULT 'pending', 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    executed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS device_groups (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

