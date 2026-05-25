package inventory

import (
	"fmt"
	"time"
)

const (
	StockAvailableKeyPrefix   = "shop:stock:available:"   // 可用库存的key前缀
	StockReservationKeyPrefix = "shop:stock:reservation:" // 预占库存的key前缀
	StockExpireQueueKey       = "shop:stock:expire_queue" // 预占库存过期队列的key, 存过期时间
	StockReservationTTL       = 30 * 24 * time.Hour
)

const (
	ReservationStatusMissing   = "MISSING"
	ReservationStatusReserved  = "RESERVED"
	ReservationStatusConfirmed = "CONFIRMED"
	ReservationStatusReleased  = "RELEASED"
)

const ReserveStockScript = `
local reservationKey = KEYS[1]
local expireQueueKey = KEYS[2]
local itemsJson = ARGV[1]
local now = ARGV[2]
local expireAt = ARGV[3]
local orderId = ARGV[4]
local items = cjson.decode(itemsJson)

local status = redis.call('HGET', reservationKey, 'status')
if status == 'RESERVED' or status == 'CONFIRMED' then
  return 1
end
if status == 'RELEASED' then
  return -3
end

for _, item in ipairs(items) do
  local stockKey = 'shop:stock:available:' .. tostring(item.sku_id)
  local current = tonumber(redis.call('GET', stockKey) or '-1')
  if current < tonumber(item.quantity) then
    return -1
  end
end

for _, item in ipairs(items) do
  local stockKey = 'shop:stock:available:' .. tostring(item.sku_id)
  redis.call('DECRBY', stockKey, tonumber(item.quantity))
end

redis.call('HSET', reservationKey,
  'status', 'RESERVED',
  'items', itemsJson,
  'created_at', now,
  'expire_at', expireAt
)
redis.call('ZADD', expireQueueKey, expireAt, orderId)
return 1
`

const ConfirmStockScript = `
local reservationKey = KEYS[1]
local expireQueueKey = KEYS[2]
local orderId = ARGV[1]
local now = ARGV[2]
local ttl = tonumber(ARGV[3])

if redis.call('EXISTS', reservationKey) == 0 then
  redis.call('ZREM', expireQueueKey, orderId)
  return 0
end
local status = redis.call('HGET', reservationKey, 'status')
if status == 'RELEASED' then
  redis.call('ZREM', expireQueueKey, orderId)
  return 2
end
if status == 'CONFIRMED' then
  redis.call('ZREM', expireQueueKey, orderId)
  return 1
end
if status ~= 'RESERVED' then
  redis.call('ZREM', expireQueueKey, orderId)
  return 0
end
redis.call('HSET', reservationKey, 'status', 'CONFIRMED', 'confirmed_at', now)
redis.call('ZREM', expireQueueKey, orderId)
redis.call('EXPIRE', reservationKey, ttl)
return 1
`

const GetReservationStatusScript = `
local reservationKey = KEYS[1]
if redis.call('EXISTS', reservationKey) == 0 then
  return 'MISSING'
end
local status = redis.call('HGET', reservationKey, 'status')
if not status then
  return 'MISSING'
end
return status
`

const ReleaseStockScript = `
local reservationKey = KEYS[1]
local expireQueueKey = KEYS[2]
local orderId = ARGV[1]
local now = ARGV[2]
local ttl = tonumber(ARGV[3])

if redis.call('EXISTS', reservationKey) == 0 then
  redis.call('ZREM', expireQueueKey, orderId)
  return 0
end

local status = redis.call('HGET', reservationKey, 'status')
if status == 'RELEASED' then
  redis.call('ZREM', expireQueueKey, orderId)
  return 1
end
if status == 'CONFIRMED' then
  redis.call('ZREM', expireQueueKey, orderId)
  return 2
end
if status ~= 'RESERVED' then
  redis.call('ZREM', expireQueueKey, orderId)
  return 0
end

local itemsJson = redis.call('HGET', reservationKey, 'items')
if not itemsJson then
  redis.call('ZREM', expireQueueKey, orderId)
  return 0
end
local items = cjson.decode(itemsJson)
for _, item in ipairs(items) do
  local stockKey = 'shop:stock:available:' .. tostring(item.sku_id)
  redis.call('INCRBY', stockKey, tonumber(item.quantity))
end
redis.call('HSET', reservationKey, 'status', 'RELEASED', 'released_at', now)
redis.call('ZREM', expireQueueKey, orderId)
redis.call('EXPIRE', reservationKey, ttl)
return 1
`

type StockItem struct {
	SkuId    uint64 `json:"sku_id" db:"sku_id"`
	Quantity uint64 `json:"quantity" db:"quantity"`
}

func AvailableStockKey(skuId uint64) string {
	return fmt.Sprintf("%s%d", StockAvailableKeyPrefix, skuId)
}

func ReservationKey(orderId uint64) string {
	return fmt.Sprintf("%s%d", StockReservationKeyPrefix, orderId)
}

func StockReservationTTLSeconds() int64 {
	return int64(StockReservationTTL / time.Second)
}
