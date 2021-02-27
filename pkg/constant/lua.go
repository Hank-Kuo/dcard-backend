package middleware

const LuaScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local numberIp = tonumber(ARGV[2])
local duration = tonumber(ARGV[3])
local hostInfo = redis.call('HGETALL', key)
local reset = tonumber(hostInfo[4])
local result = {}

if #hostInfo == 0 or reset < now then
    reset = now + duration
    redis.call('HMSET', key, "count", 1, "reset", reset)
    result[1] = numberIp - 1
    result[2] = reset
    return result
end

local count = tonumber(hostInfo[2])

if count < numberIp then
    local newCount = redis.call('HINCRBY', key, "count", 1)	
    result[1] = numberIp - newCount
    result[2] = reset
    return result
else
    result[1] = -1
    result[2] = reset
    return result
end
`
