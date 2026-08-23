-- turtle_client.lua
-- Bridge client for Robot-Project bot server.
-- Registers itself once, then loops: report status -> fetch command -> execute -> sleep.

----------------------------------------------------------------------
-- CONFIG
----------------------------------------------------------------------

local BASE_URL      = "http://localhost:8080"  -- <-- set to your server's LAN IP
local ID_FILE       = ".bot_id"                -- persisted on the turtle's own disk
local POLL_INTERVAL = 1                        -- seconds between loop iterations
local LOW_FUEL_WARN = 50                       -- log a warning below this fuel level

----------------------------------------------------------------------
-- HTTP HELPERS
-- CC:Tweaked's http.get / http.post are blocking convenience wrappers,
-- which suits a simple poll loop like this one.
----------------------------------------------------------------------

-- CC:Tweaked's http.get/http.post return nil on any non-2xx response, with
-- `err` set to just the reason phrase ("Bad Request", "Not Found", etc).
-- They also return a THIRD value on failure: the failed response handle,
-- which still has a readable status code and body. We read that here so
-- we can see the server's actual JSON error message instead of guessing.
local function readFailedResponse(failRes)
    if not failRes then return nil, nil end
    local code = failRes.getResponseCode()
    local body = failRes.readAll()
    failRes.close()
    return code, body
end

local function httpGet(path)
    local res, err, failRes = http.get(BASE_URL .. path)
    if not res then
        local code, body = readFailedResponse(failRes)
        print("[http] GET " .. path .. " failed: " .. tostring(err) ..
              (code and (" (HTTP " .. code .. ")") or "") ..
              (body and body ~= "" and (" body: " .. body) or ""))
        return nil, code, body
    end
    local body = res.readAll()
    res.close()
    return textutils.unserializeJSON(body), res.getResponseCode(), body
end

-- `bodyTable` can be a Lua table (auto-serialized) or, when rawJson is true,
-- an already-built JSON string. The raw-string path exists because
-- textutils.serializeJSON can't reliably tell an empty array apart from an
-- empty object, so payloads with array fields that might be empty (like
-- our inventory list) are built manually instead. See reportStatus().
local function httpPost(path, bodyTable, rawJson)
    local payload
    if rawJson then
        payload = bodyTable or ""
    else
        payload = bodyTable and textutils.serializeJSON(bodyTable) or ""
    end
    local headers = (bodyTable and payload ~= "") and { ["Content-Type"] = "application/json" } or nil

    local res, err, failRes = http.post(BASE_URL .. path, payload, headers)
    if not res then
        local code, body = readFailedResponse(failRes)
        print("[http] POST " .. path .. " failed: " .. tostring(err) ..
              (code and (" (HTTP " .. code .. ")") or "") ..
              (body and body ~= "" and (" body: " .. body) or ""))
        return code, nil, body
    end
    local code = res.getResponseCode()
    local body = res.readAll()
    res.close()
    local decoded = nil
    if body and #body > 0 then
        decoded = textutils.unserializeJSON(body)
    end
    return code, decoded, body
end

----------------------------------------------------------------------
-- "BOT NOT FOUND" DETECTION
-- The server can forget a bot ID (e.g. it was restarted, or its bot list
-- was cleared) while the turtle still has that ID cached on disk. When
-- that happens every request 404s forever unless we notice and
-- re-register. We check the decoded/raw error body for the server's
-- "Bot not found" message rather than trusting a bare 404 status code,
-- since a 404 could in principle mean something else (e.g. a bad path).
----------------------------------------------------------------------

local function isBotNotFound(code, rawBody)
    if code ~= 404 or not rawBody then return false end
    return rawBody:find("Bot not found", 1, true) ~= nil
end

----------------------------------------------------------------------
-- ID REGISTRATION
-- Reuses a saved ID across reboots instead of re-registering every time.
----------------------------------------------------------------------

local function loadSavedId()
    if fs.exists(ID_FILE) then
        local f = fs.open(ID_FILE, "r")
        local id = f.readAll()
        f.close()
        if id and #id > 0 then
            return id
        end
    end
    return nil
end

local function saveId(id)
    local f = fs.open(ID_FILE, "w")
    f.write(id)
    f.close()
end

local function clearSavedId()
    if fs.exists(ID_FILE) then
        fs.delete(ID_FILE)
    end
end

local function registerNewId()
    local idResp = httpGet("/id")
    if not idResp or not idResp.id then
        error("Could not get a new ID from server")
    end

    local code = httpPost("/id/" .. idResp.id)
    if code ~= 200 and code ~= 201 then
        error("Registration failed, server returned code " .. tostring(code))
    end

    saveId(idResp.id)
    return idResp.id
end

local function getBotId()
    local id = loadSavedId()
    if id then
        print("[init] Using saved ID: " .. id)
        return id
    end
    id = registerNewId()
    print("[init] Registered new ID: " .. id)
    return id
end

-- Wipes the stale saved ID and registers a fresh one. Called whenever the
-- server tells us it no longer knows about our current bot ID.
local function reregister()
    print("[init] Server doesn't recognize our bot ID, re-registering...")
    clearSavedId()
    local id = registerNewId()
    print("[init] Registered new ID: " .. id)
    return id
end

----------------------------------------------------------------------
-- STATUS REPORTING
----------------------------------------------------------------------

-- textutils.serializeJSON can't distinguish an empty Lua table from an
-- empty JSON object, so an empty inventory always serializes as `{}` no
-- matter what metatable hints are set. The server wants a JSON array
-- ([]bot.InventoryItem in Go), so `{}` triggers:
--   "json: cannot unmarshal object into Go struct field Status.inventory
--    of type []bot.InventoryItem"
-- Workaround: build the inventory portion of the JSON as a raw string,
-- letting serializeJSON handle just the individual item objects (which
-- are never ambiguous since they always have named keys).
local function collectInventory()
    local inventory = {}
    for slot = 1, 16 do
        local detail = turtle.getItemDetail(slot)
        if detail then
            table.insert(inventory, {
                slot  = slot,
                item  = detail.name,  -- server's bot.InventoryItem expects "item", not "name"
                count = detail.count,
            })
        end
    end
    return inventory
end

local function inventoryToJsonArray(inventory)
    if #inventory == 0 then
        return "[]"
    end
    local parts = {}
    for _, item in ipairs(inventory) do
        table.insert(parts, textutils.serializeJSON(item))
    end
    return "[" .. table.concat(parts, ",") .. "]"
end

-- Returns the HTTP status code (and, on a "bot not found" error, a second
-- return value of true) so the caller can decide whether to re-register.
local function reportStatus(botId)
    local fuel = turtle.getFuelLevel()
    local fuelLimit = turtle.getFuelLimit()

    -- Map "unlimited" fuel to 0 to prevent 400 validation errors on backend integers
    if fuel == "unlimited" then fuel = 0 end
    if fuelLimit == "unlimited" then fuelLimit = 0 end

    if type(fuel) == "number" and fuel < LOW_FUEL_WARN and fuel > 0 then
        print("[status] WARNING: low fuel (" .. fuel .. ")")
    end

    local inventoryJson = inventoryToJsonArray(collectInventory())
    local payload = string.format(
        '{"fuel":%d,"fuel_limit":%d,"inventory":%s}',
        fuel, fuelLimit, inventoryJson
    )

    local code, _, rawBody = httpPost("/id/" .. botId .. "/status", payload, true)

    -- Surface non-2xx responses instead of failing silently.
    if code and (code < 200 or code >= 300) then
        print("[status] WARNING: server returned " .. tostring(code))
    end

    return code, isBotNotFound(code, rawBody)
end

----------------------------------------------------------------------
-- COMMAND DISPATCH TABLE
-- Maps a "payload" string (sent by the controller) to a turtle action.
----------------------------------------------------------------------

local commands = {
    -- movement
    forward   = function() return turtle.forward() end,
    back      = function() return turtle.back() end,
    up        = function() return turtle.up() end,
    down      = function() return turtle.down() end,
    turnLeft  = function() return turtle.turnLeft() end,
    turnRight = function() return turtle.turnRight() end,

    -- digging
    dig       = function() return turtle.dig() end,
    digUp     = function() return turtle.digUp() end,
    digDown   = function() return turtle.digDown() end,

    -- placing
    place     = function() return turtle.place() end,
    placeUp   = function() return turtle.placeUp() end,
    placeDown = function() return turtle.placeDown() end,

    -- attacking
    attack     = function() return turtle.attack() end,
    attackUp   = function() return turtle.attackUp() end,
    attackDown = function() return turtle.attackDown() end,

    -- detection (block presence)
    detect     = function() return turtle.detect() end,
    detectUp   = function() return turtle.detectUp() end,
    detectDown = function() return turtle.detectDown() end,

    -- inspection (block data)
    inspect     = function() return turtle.inspect() end,
    inspectUp   = function() return turtle.inspectUp() end,
    inspectDown = function() return turtle.inspectDown() end,

    -- comparison (selected slot vs world block)
    compare     = function() return turtle.compare() end,
    compareUp   = function() return turtle.compareUp() end,
    compareDown = function() return turtle.compareDown() end,

    -- sucking items from inventories/ground (with safe number conversion)
    suck     = function(arg) return turtle.suck(tonumber(arg)) end,
    suckUp   = function(arg) return turtle.suckUp(tonumber(arg)) end,
    suckDown = function(arg) return turtle.suckDown(tonumber(arg)) end,

    -- dropping items
    drop     = function(arg) return turtle.drop(tonumber(arg)) end,
    dropUp   = function(arg) return turtle.dropUp(tonumber(arg)) end,
    dropDown = function(arg) return turtle.dropDown(tonumber(arg)) end,

    -- refueling
    refuel = function(arg) return turtle.refuel(tonumber(arg)) end,

    -- crafting (crafty turtles only)
    craft = function(arg) return turtle.craft(tonumber(arg)) end,

    -- equipping tools/peripherals in each hand
    equipLeft  = function() return turtle.equipLeft() end,
    equipRight = function() return turtle.equipRight() end,

    -- inventory slot management
    select = function(arg)
        local slot = tonumber(arg)
        if not slot then return false, "select requires a numeric slot" end
        return turtle.select(slot)
    end,
    transferTo = function(arg)
        local slot = tonumber(arg)
        if not slot then return false, "transferTo requires a numeric slot" end
        return turtle.transferTo(slot)
    end,

    -- no-op, useful as an explicit "do nothing this cycle" command
    stop = function() return true end,
}

----------------------------------------------------------------------
-- COMMAND EXECUTION
----------------------------------------------------------------------

local lastExecutedAt = nil

local function executeCommand(msg)
    if not msg then return end

    -- msg.timestamp lets us avoid re-running the same command every poll
    if msg.timestamp == lastExecutedAt then return end
    lastExecutedAt = msg.timestamp

    local action = commands[msg.payload]
    if not action then
        print("[cmd] Unknown command: " .. tostring(msg.payload))
        return
    end

    print("[cmd] Executing: " .. msg.payload)
    local ok, errOrResult = action(msg.arg)
    if ok then
        print("[cmd] OK")
    else
        print("[cmd] Failed: " .. tostring(errOrResult))
    end
end

----------------------------------------------------------------------
-- MAIN LOOP
----------------------------------------------------------------------

local function main()
    local botId = getBotId()
    print("[init] Bot online as " .. botId)

    while true do
        local _, statusNotFound = reportStatus(botId)

        if statusNotFound then
            -- No point polling for a command with an ID the server has
            -- already forgotten; re-register now and pick up next loop.
            botId = reregister()
        else
            local resp, msgCode, rawBody = httpGet("/id/" .. botId .. "/message")
            if isBotNotFound(msgCode, rawBody) then
                botId = reregister()
            elseif resp and resp.message then
                executeCommand(resp.message)
            end
        end

        sleep(POLL_INTERVAL)
    end
end

main()