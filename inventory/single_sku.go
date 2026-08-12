package inventory

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Single-SKU inventory is the shared Web3 contract. The legacy multi-SKU
// scripts in inventory.go remain unchanged for existing consumers.
const (
	AllocationStateReserved = "A"
	AllocationStateReleased = "R"
	AllocationStateConsumed = "C"
	StockPressureQueueKey   = "inv:pressure:sku_queue"
)

const (
	ResultInsufficient int64 = 0
	ResultApplied      int64 = 1
	ResultIdempotent   int64 = 2
	ResultConflict     int64 = -2
	ResultInvalidState int64 = -3
	ResultDataMissing  int64 = -4
	ResultInvalidInput int64 = -5
)

// Durable allocation states are shared by the API producer and event worker.
const (
	AllocationSagaStateAllocatePending uint8 = 0
	AllocationSagaStateAllocated       uint8 = 10
	AllocationSagaStateReleasePending  uint8 = 20
	AllocationSagaStateReleased        uint8 = 30
	AllocationSagaStateConsumePending  uint8 = 40
	AllocationSagaStateConsumed        uint8 = 50
	AllocationSagaStateFailed          uint8 = 90

	AllocationPurposeSale uint8 = 1
)

// AllocateOneScript atomically reserves one SKU. Redis keeps only a temporary
// idempotency field; MySQL inventory_allocation owns the durable lifecycle.
const AllocateOneScript = `
local qty = tonumber(ARGV[1])
local orderId = ARGV[2]
if not qty or qty <= 0 or not orderId or orderId == '' then return -5 end
local value = redis.call('HGET', KEYS[3], orderId)
if value then
  local sep = string.find(value, ':', 1, true)
  if not sep then return -4 end
  local state = string.sub(value, 1, sep - 1)
  local storedQty = tonumber(string.sub(value, sep + 1))
  if not storedQty then return -4 end
  if storedQty ~= qty then return -2 end
  if state == 'A' then return 2 end
  return -3
end
local available = tonumber(redis.call('GET', KEYS[1]))
if not available then return -4 end
if available < qty then return 0 end
redis.call('DECRBY', KEYS[1], qty)
redis.call('INCRBY', KEYS[2], qty)
redis.call('HSET', KEYS[3], orderId, 'A:' .. tostring(qty))
return 1
`

// ReleaseOneScript keeps R:qty only across the Redis-success/MySQL-terminal
// crash window. The caller deletes the field after MySQL commits RELEASED.
const ReleaseOneScript = `
local orderId = ARGV[1]
if not orderId or orderId == '' then return -5 end
local value = redis.call('HGET', KEYS[3], orderId)
if not value then return 0 end
local sep = string.find(value, ':', 1, true)
if not sep then return -4 end
local state = string.sub(value, 1, sep - 1)
local qty = tonumber(string.sub(value, sep + 1))
if not qty or qty <= 0 then return -4 end
if state == 'R' then return 2 end
if state ~= 'A' then return -3 end
local available = tonumber(redis.call('GET', KEYS[1]))
local reserved = tonumber(redis.call('GET', KEYS[2]))
if not available or not reserved or reserved < qty then return -4 end
redis.call('INCRBY', KEYS[1], qty)
redis.call('DECRBY', KEYS[2], qty)
redis.call('HSET', KEYS[3], orderId, 'R:' .. tostring(qty))
return 1
`

// ConsumeOneScript keeps C:qty only across the Redis-success/MySQL-terminal
// crash window. The caller deletes the field after MySQL commits CONSUMED.
const ConsumeOneScript = `
local orderId = ARGV[1]
if not orderId or orderId == '' then return -5 end
local value = redis.call('HGET', KEYS[2], orderId)
if not value then return 0 end
local sep = string.find(value, ':', 1, true)
if not sep then return -4 end
local state = string.sub(value, 1, sep - 1)
local qty = tonumber(string.sub(value, sep + 1))
if not qty or qty <= 0 then return -4 end
if state == 'C' then return 2 end
if state ~= 'A' then return -3 end
local reserved = tonumber(redis.call('GET', KEYS[1]))
if not reserved or reserved < qty then return -4 end
redis.call('DECRBY', KEYS[1], qty)
redis.call('HSET', KEYS[2], orderId, 'C:' .. tostring(qty))
return 1
`

// CleanupOneScript removes an idempotency field only when it exactly matches
// the terminal value committed by MySQL. A missing field is already clean.
const CleanupOneScript = `
local value = redis.call('HGET', KEYS[1], ARGV[1])
if not value then return 2 end
if value ~= ARGV[2] then return -2 end
redis.call('HDEL', KEYS[1], ARGV[1])
return 1
`

// UnlockOrderScript and RenewOrderScript protect owner-based paying guards.
const UnlockOrderScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('DEL', KEYS[1])
end
return 0
`

const RenewOrderScript = `
if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('EXPIRE', KEYS[1], ARGV[2])
end
return 0
`

var (
	AllocateOneSHA = scriptSHA(AllocateOneScript)
	ReleaseOneSHA  = scriptSHA(ReleaseOneScript)
	ConsumeOneSHA  = scriptSHA(ConsumeOneScript)
	CleanupOneSHA  = scriptSHA(CleanupOneScript)
	UnlockOrderSHA = scriptSHA(UnlockOrderScript)
	RenewOrderSHA  = scriptSHA(RenewOrderScript)
)

func AvailableKey(skuID uint64) string {
	return fmt.Sprintf("inv:{%d}:available", skuID)
}

func ReservedKey(skuID uint64) string {
	return fmt.Sprintf("inv:{%d}:reserved", skuID)
}

func AllocationsKey(skuID uint64) string {
	return fmt.Sprintf("inv:{%d}:allocations", skuID)
}

func AllocationField(orderID uint64) string {
	return strconv.FormatUint(orderID, 10)
}

func AllocationValue(state string, quantity uint64) string {
	return state + ":" + strconv.FormatUint(quantity, 10)
}

func ParseAllocationValue(value string) (state string, quantity uint64, ok bool) {
	state, rawQuantity, found := strings.Cut(value, ":")
	if !found || (state != AllocationStateReserved && state != AllocationStateReleased && state != AllocationStateConsumed) {
		return "", 0, false
	}
	quantity, err := strconv.ParseUint(rawQuantity, 10, 64)
	if err != nil || quantity == 0 {
		return "", 0, false
	}
	return state, quantity, true
}

func scriptSHA(script string) string {
	digest := sha1.Sum([]byte(script))
	return hex.EncodeToString(digest[:])
}
