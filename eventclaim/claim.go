package eventclaim

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	markerPrefix         = "claim:"
	defaultLeaseDuration = 2 * time.Minute
	minimumLeaseDuration = 30 * time.Second
)

var (
	ErrClaimLost       = errors.New("event claim lost")
	ErrHeartbeatFailed = errors.New("event claim heartbeat failed")
	ErrPermanent       = errors.New("event handling failed permanently")
)

type permanentError struct{ err error }

func (e *permanentError) Error() string { return e.err.Error() }
func (e *permanentError) Unwrap() error { return e.err }
func (e *permanentError) Is(target error) bool {
	return target == ErrPermanent
}

func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return &permanentError{err: err}
}

func Permanentf(format string, args ...any) error {
	return Permanent(fmt.Errorf(format, args...))
}

type markerContextKey struct{}

func WithNewMarker(ctx context.Context) (context.Context, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("generate event claim marker: %w", err)
	}
	return context.WithValue(ctx, markerContextKey{}, markerPrefix+hex.EncodeToString(buf)), nil
}

func Marker(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("event claim context is nil")
	}
	marker, _ := ctx.Value(markerContextKey{}).(string)
	if !strings.HasPrefix(marker, markerPrefix) || len(marker) == len(markerPrefix) {
		return "", errors.New("event claim marker is missing")
	}
	return marker, nil
}

func ParseLeaseDuration(raw string) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultLeaseDuration
	}
	duration, err := time.ParseDuration(raw)
	if err != nil || duration < minimumLeaseDuration {
		return defaultLeaseDuration
	}
	return duration
}

func HeartbeatInterval(leaseDuration time.Duration) time.Duration {
	interval := leaseDuration / 4
	if interval < 5*time.Second {
		return 5 * time.Second
	}
	if interval > 30*time.Second {
		return 30 * time.Second
	}
	return interval
}

func LeaseSeconds(leaseDuration time.Duration) int64 {
	return int64((leaseDuration + time.Second - 1) / time.Second)
}

func HeartbeatTimeout(leaseDuration time.Duration) time.Duration {
	timeout := HeartbeatInterval(leaseDuration)
	if timeout > 10*time.Second {
		return 10 * time.Second
	}
	return timeout
}

func LostError(eventID uint64) error {
	return fmt.Errorf("%w: event_id=%d", ErrClaimLost, eventID)
}
