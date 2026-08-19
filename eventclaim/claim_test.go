package eventclaim

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWithNewMarker(t *testing.T) {
	ctx, err := WithNewMarker(context.Background())
	if err != nil {
		t.Fatalf("WithNewMarker() error = %v", err)
	}
	marker, err := Marker(ctx)
	if err != nil {
		t.Fatalf("Marker() error = %v", err)
	}
	if !strings.HasPrefix(marker, markerPrefix) {
		t.Fatalf("Marker() = %q", marker)
	}
	if _, err = Marker(context.Background()); err == nil {
		t.Fatal("Marker() expected missing marker error")
	}
	secondCtx, err := WithNewMarker(context.Background())
	if err != nil {
		t.Fatalf("second WithNewMarker() error = %v", err)
	}
	secondMarker, err := Marker(secondCtx)
	if err != nil {
		t.Fatalf("second Marker() error = %v", err)
	}
	if marker == secondMarker {
		t.Fatal("WithNewMarker() returned duplicate markers")
	}
}

func TestParseLeaseDuration(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want time.Duration
	}{
		{name: "default", want: defaultLeaseDuration},
		{name: "valid", raw: "45s", want: 45 * time.Second},
		{name: "below minimum", raw: "5s", want: defaultLeaseDuration},
		{name: "invalid", raw: "invalid", want: defaultLeaseDuration},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ParseLeaseDuration(test.raw); got != test.want {
				t.Fatalf("ParseLeaseDuration(%q) = %s, want %s", test.raw, got, test.want)
			}
		})
	}
}

func TestHeartbeatInterval(t *testing.T) {
	if got := HeartbeatInterval(40 * time.Second); got != 10*time.Second {
		t.Fatalf("HeartbeatInterval() = %s", got)
	}
	if got := HeartbeatInterval(10 * time.Minute); got != 30*time.Second {
		t.Fatalf("HeartbeatInterval() capped = %s", got)
	}
}

func TestLeaseSecondsRoundsUp(t *testing.T) {
	if got := LeaseSeconds(30500 * time.Millisecond); got != 31 {
		t.Fatalf("LeaseSeconds() = %d", got)
	}
}

func TestHeartbeatTimeout(t *testing.T) {
	if got := HeartbeatTimeout(40 * time.Second); got != 10*time.Second {
		t.Fatalf("HeartbeatTimeout() = %s", got)
	}
	if got := HeartbeatTimeout(30 * time.Second); got != 7500*time.Millisecond {
		t.Fatalf("HeartbeatTimeout() short lease = %s", got)
	}
}

func TestLostError(t *testing.T) {
	if err := LostError(42); !errors.Is(err, ErrClaimLost) {
		t.Fatalf("LostError() = %v", err)
	}
}

func TestPermanent(t *testing.T) {
	cause := errors.New("invalid projection")
	err := Permanent(cause)
	if !errors.Is(err, ErrPermanent) {
		t.Fatalf("Permanent() = %v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("Permanent() lost cause: %v", err)
	}
	if err.Error() != "invalid projection" {
		t.Fatalf("Permanent() error = %q", err.Error())
	}
	if Permanent(nil) != nil {
		t.Fatal("Permanent(nil) expected nil")
	}
}
