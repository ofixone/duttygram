#!/usr/bin/env tarantool

local function getEnv(key, default)
    local value = os.getenv(key)
    if value == nil then
        return default
    end

    return value
end

local config = {
    memory_limit = tonumber(getEnv("TNT_MEMORY_LIMIT", 512 * 1024 * 1024)),
    listen = tonumber(getEnv("TNT_LISTEN", 3301)),
    log_level = tonumber(getEnv("TNT_LOG_LEVEL", 2)),
    log_format = getEnv("TNT_LOG_FORMAT", "json"),
    user = getEnv("TNT_USER", "tarantool"),
    password = getEnv("TNT_PASSWORD", "tarantool"),
    migration_user = getEnv("TNT_MIGRATION_USER", "migrator"),
    migration_password = getEnv("TNT_MIGRATION_PASSWORD", "migrator"),
}

box.cfg {
    listen = config.listen,
    memtx_memory = config.memory_limit,
    read_only = config.is_readonly,
    log_level = config.log_level,
    log_format = config.log_format,
    wal_dir = "/var/lib/tarantool",
    memtx_dir = "/var/lib/tarantool",
    vinyl_dir = "/var/lib/tarantool",
}

box.once("bootstrap", function()
    -- Основной пользователь
    box.schema.user.create(config.user, { password = config.password })
    box.schema.user.grant(config.user, "read,write,execute", "universe")
end)

box.once("migration", function()
    -- Пользователь для миграций
    box.schema.user.create(config.migration_user, { password = config.migration_password })
    box.schema.user.grant(config.migration_user, "read,write,create,alter,drop,execute", "universe")
end)

box.once("backup", function()
    box.schema.func.create("backup_start")
    box.schema.func.create("backup_stop")

    -- Пользователь для бэкапов
    box.schema.user.create("backup", { password = "backup" })
    box.schema.user.grant("backup", "execute", "function", "backup_start")
    box.schema.user.grant("backup", "execute", "function", "backup_stop")
end)

function backup_start()
    return box.backup.start()
end

function backup_stop()
    return box.backup.stop()
end

API = require("api")
migrator = require("migrations/migrator")

local run_migrations = os.getenv("TNT_RUN_MIGRATIONS")
if run_migrations == "true" then
    migrator.migrate(require("migrations/migrations"))
end