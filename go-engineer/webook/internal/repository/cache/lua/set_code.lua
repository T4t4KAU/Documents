-- 发送到的 key，也就是 code:业务:手机号码
local key = KEYS[1]
-- 使用次数，也就是验证次数
local cntKey = key..":cnt"
local val = ARGV[1]
-- 验证码的有效时间是十分钟，600 秒
local ttl = tonumber(redis.call("ttl", key))

-- -1 是 key 存在，但是没有过期时间
if ttl == -1 then
    -- 有人误操作，导致 key 冲突
    return -2
-- -2 是 key 不存在，ttl < 540 是发了一个验证码，已经超过一分钟了，可以重新发送
elseif ttl == -2 or ttl < 540 then
    --     后续如果验证码有不同的过期时间，要在这里优化
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 已经发送了一个验证码，但是还不到一分钟
    return -1
end