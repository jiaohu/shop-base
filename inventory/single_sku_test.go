package inventory

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"testing"
)

func TestSingleSKUContract(t *testing.T) {
	if got := AvailableKey(42); got != "inv:{42}:available" {
		t.Fatalf("unexpected available key: %s", got)
	}
	if got := ReservedKey(42); got != "inv:{42}:reserved" {
		t.Fatalf("unexpected reserved key: %s", got)
	}
	if got := AllocationsKey(42); got != "inv:{42}:allocations" {
		t.Fatalf("unexpected allocations key: %s", got)
	}
	if got := AllocationField(7); got != "7" {
		t.Fatalf("unexpected allocation field: %s", got)
	}
	if got := AllocationValue(AllocationStateConsumed, 17); got != "C:17" {
		t.Fatalf("unexpected allocation value: %s", got)
	}
	state, quantity, ok := ParseAllocationValue("A:17")
	if !ok || state != AllocationStateReserved || quantity != 17 {
		t.Fatalf("unexpected parsed allocation: state=%s quantity=%d ok=%v", state, quantity, ok)
	}
	if _, _, ok = ParseAllocationValue("broken"); ok {
		t.Fatal("broken allocation value must be rejected")
	}

	for name, script := range map[string]string{
		"allocate": AllocateOneScript,
		"release":  ReleaseOneScript,
		"consume":  ConsumeOneScript,
		"cleanup":  CleanupOneScript,
	} {
		t.Run(name, func(t *testing.T) {
			lower := strings.ToLower(script)
			for _, forbidden := range []string{"cjson", "for ", "while ", "zadd", "zrem", "zrange", "scan", "expire"} {
				if strings.Contains(lower, forbidden) {
					t.Fatalf("%s script contains forbidden operation %q", name, forbidden)
				}
			}
			if lines := strings.Count(strings.TrimSpace(script), "\n") + 1; lines > 32 {
				t.Fatalf("%s script is too large: %d lines", name, lines)
			}
		})
	}

	assertScriptSHA(t, AllocateOneScript, AllocateOneSHA)
	assertScriptSHA(t, ReleaseOneScript, ReleaseOneSHA)
	assertScriptSHA(t, ConsumeOneScript, ConsumeOneSHA)
	assertScriptSHA(t, CleanupOneScript, CleanupOneSHA)
	assertScriptSHA(t, UnlockOrderScript, UnlockOrderSHA)
	assertScriptSHA(t, RenewOrderScript, RenewOrderSHA)
}

func TestAllocationSagaStateValues(t *testing.T) {
	values := []uint8{
		AllocationSagaStateAllocatePending,
		AllocationSagaStateAllocated,
		AllocationSagaStateReleasePending,
		AllocationSagaStateReleased,
		AllocationSagaStateConsumePending,
		AllocationSagaStateConsumed,
		AllocationSagaStateFailed,
		AllocationPurposeSale,
	}
	want := []uint8{0, 10, 20, 30, 40, 50, 90, 1}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("state %d: want %d, got %d", i, want[i], values[i])
		}
	}
}

func assertScriptSHA(t *testing.T, script, want string) {
	t.Helper()
	digest := sha1.Sum([]byte(script))
	if got := hex.EncodeToString(digest[:]); got != want {
		t.Fatalf("script SHA mismatch: want %s, got %s", want, got)
	}
}
