CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT
);

-- статус: pending, sent, executed, error, online, offline
CREATE TABLE IF NOT EXISTS statuses (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT,
    code TEXT NOT NULL UNIQUE,
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
    -- user_id TEXT REFERENCES users(id), ВОПРОС НУЖНО ЛИ
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
    installed_programs JSONB, -- Новый столбец для списка установленных программ
    cpu_percent REAL,
    bytes_sent BIGINT,
    bytes_recv BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    -- В будущем можно реализовать партиционирование по created_at для таймсерийных данных
);

-- Таблица команд, отправленных устройствам
CREATE TABLE IF NOT EXISTS commands (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    command_type TEXT NOT NULL,
    parameters JSONB,
    initiator TEXT REFERENCES users(id),
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