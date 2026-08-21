-- turtle_client.lua
-- Bridge client for Robot-Project bot server.
-- Registers itself once, then loops: report status -> fetch command -> execute -> sleep.

----------------------------------------------------------------------
-- CONFIG
----------------------------------------------------------------------

local BASE_URL     = "http://localhost:8080"  -- <-- set to your server's LAN IP
local ID_FILE      = ".bot_id"                    -- persisted on the turtle's own disk
local POLL_INTERVAL = 1                           -- seconds between loop iterations
local LOW_FUEL_WARN = 50                          -- log a warning below this fuel level

----------------------------------------------------------------------
-- HTTP HELPERS
-- CC:Tweaked's http.get / http.post are blocking convenience wrappers
-- (as opposed to the async http.request + "http_success" event API),
-- which suits a simple poll loop like this one.
----------------------------------------------------------------------

local function httpGet(path)
    local res, err = http.get(BASE_URL .. path)
    if not res then
        print("[http] GET " .. path .. " failed: " .. tostring(err))
        return nil
    end
    local body = res.readAll()
    res.close()
    return textutils.unserializeJSON(body)
end

local function httpPost(path, bodyTable)
    local payload = bodyTable and textutils.serializeJSON(bodyTable) or nil
    local res, err = http.post(
        BASE_URL .. path,
        payload,
        { ["Content-Type"] = "application/json" }
    )
    if not res then
        print("[http] POST " .. path .. " failed: " .. tostring(err))
        return nil, nil
    end
    local code = res.getResponseCode()
    local body = res.readAll()
    res.close()
    local decoded = nil
    if body and #body > 0 then
        decoded = textutils.unserializeJSON(body)
    end
    return code, decoded
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

----------------------------------------------------------------------
-- STATUS REPORTING
----------------------------------------------------------------------

local function collectInventory()
    local inventory = {}
    for slot = 1, 16 do
        local detail = turtle.getItemDetail(slot)
        if detail then
            table.insert(inventory, {
                slot  = slot,
                name  = detail.name,
                count = detail.count,
            })
        end
    end
    return inventory
end

local function reportStatus(botId)
    local fuel = turtle.getFuelLevel()
    local fuelLimit = turtle.getFuelLimit()

    -- unlimited fuel worlds return the string "unlimited" instead of a number
    if fuel == "unlimited" then fuel = -1 end
    if fuelLimit == "unlimited" then fuelLimit = -1 end

    if type(fuel) == "number" and fuel < LOW_FUEL_WARN and fuel >= 0 then
        print("[status] WARNING: low fuel (" .. fuel .. ")")
    end

    httpPost("/id/" .. botId .. "/status", {
        fuel       = fuel,
        fuel_limit = fuelLimit,
        inventory  = collectInventory(),
    })
end

----------------------------------------------------------------------
-- COMMAND DISPATCH TABLE
-- Maps a "payload" string (sent by the controller) to a turtle action.
-- Every entry returns (ok, errMsg) matching the turtle API's own
-- return signature where applicable, so results can be logged uniformly.
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

    -- sucking items from inventories/ground
    suck     = function(arg) return turtle.suck(arg) end,
    suckUp   = function(arg) return turtle.suckUp(arg) end,
    suckDown = function(arg) return turtle.suckDown(arg) end,

    -- dropping items
    drop     = function(arg) return turtle.drop(arg) end,
    dropUp   = function(arg) return turtle.dropUp(arg) end,
    dropDown = function(arg) return turtle.dropDown(arg) end,

    -- refueling
    refuel = function(arg) return turtle.refuel(arg) end,

    -- crafting (crafty turtles only)
    craft = function(arg) return turtle.craft(arg) end,

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
        reportStatus(botId)

        local resp = httpGet("/id/" .. botId .. "/message")
        if resp and resp.message then
            executeCommand(resp.message)
        end

        sleep(POLL_INTERVAL)
    end
end

main()