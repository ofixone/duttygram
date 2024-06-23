--- Запускатель миграций
local migrator = {}

--- Запускает набор миграций
--- @param migrations table<table<name string, func fun()>> Миграции
function migrator.migrate(migrations)
    --- @type log
    local log = require('log')
    log.info("[migrator] start migrations")
    local total_migrations = 0
    local index = 1
    for _, migration_obj in ipairs(migrations) do
        local name = migration_obj["name"]
        local func = migration_obj["func"]

        assert(type(name) == "string", "migration name should be string")
        assert(type(func) == "function", "func should be a function")

        local migration_name = name .. "_" .. tostring(index)
        index = index + 1
        box.once(migration_name, function()
            total_migrations = total_migrations + 1
            log.info("[migrator] start migrate '%s'", migration_name)
            local status, err = pcall(func)
            if not status then
                log.error("[migrator] error while migrating, database is in inconsistent state. This migration" ..
                        " has been reset and cancel whole migration process: %s", err)
                box.space._schema:delete("once" .. migration_name)
                box.error(err)
            end
            log.info("[migrator] finish migrate '%s", migration_name)
        end)
    end
    log.info("[migrator] finish migrations, migrates have ran %d", total_migrations)
end

return migrator