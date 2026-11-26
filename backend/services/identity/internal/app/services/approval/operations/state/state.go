// Package state owns pure workflow transitions. It does not authenticate
// reviewers, persist decisions, or apply resource mutations.
package state

import (
	"errors"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"grovehub/services/identity/internal/app/constants/approval"
)

var (
	ErrInvalid = errors.New("invalid approval transition")
	ErrVersion = errors.New("approval version conflict")
	ErrForbidden = errors.New("approval actor forbidden")
	ErrExpired = errors.New("approval request expired")
)

type Snapshot struct {
	TenantID string
	RequesterID string
	Status approval.Status
	CurrentStep int
	StepCount int
	Version int64
	ExpiresAt time.Time
}

type Actor struct {
	TenantID string
	AccountID string
}

func validate(current Snapshot, version int64, now time.Time) error {
	if current.TenantID == "" || current.RequesterID == "" || current.Version < 1 || current.Version == math.MaxInt64 || current.StepCount < 1 || current.StepCount > approval.MaxSteps || current.CurrentStep < 1 || current.CurrentStep > current.StepCount || now.IsZero() || current.ExpiresAt.IsZero() {
		return ErrInvalid
	}
	if version != current.Version {
		return ErrVersion
	}
	return nil
}

func Decide(current Snapshot, actor Actor, version int64, decision approval.Decision, comment string, now time.Time) (Snapshot, error) {
	if err := validate(current, version, now); err != nil {
		return Snapshot{}, err
	}
	if actor.TenantID != current.TenantID || actor.AccountID == "" || actor.AccountID == current.RequesterID {
		return Snapshot{}, ErrForbidden
	}
	if current.Status != approval.Pending {
		return Snapshot{}, ErrInvalid
	}
	if !now.Before(current.ExpiresAt) {
		return Snapshot{}, ErrExpired
	}
	if !utf8.ValidString(comment) || strings.ContainsRune(comment, 0) || utf8.RuneCountInString(comment) > 2000 {
		return Snapshot{}, ErrInvalid
	}
	switch decision {
	case approval.Reject:
		if strings.TrimSpace(comment) == "" {
			return Snapshot{}, ErrInvalid
		}
		current.Status = approval.Rejected
	case approval.Approve:
		if current.CurrentStep == current.StepCount {
			current.Status = approval.Approved
		} else {
			current.CurrentStep++
		}
	default:
		return Snapshot{}, ErrInvalid
	}
	current.Version++
	return current, nil
}

func Cancel(current Snapshot, actor Actor, version int64, now time.Time) (Snapshot, error) {
	if err := validate(current, version, now); err != nil {
		return Snapshot{}, err
	}
	if actor.TenantID != current.TenantID || actor.AccountID != current.RequesterID {
		return Snapshot{}, ErrForbidden
	}
	if current.Status != approval.Pending {
		return Snapshot{}, ErrInvalid
	}
	if !now.Before(current.ExpiresAt) {
		return Snapshot{}, ErrExpired
	}
	current.Status = approval.Cancelled
	current.Version++
	return current, nil
}

func Expire(current Snapshot, version int64, now time.Time) (Snapshot, error) {
	if err := validate(current, version, now); err != nil {
		return Snapshot{}, err
	}
	if (current.Status != approval.Pending && current.Status != approval.Approved) || now.Before(current.ExpiresAt) {
		return Snapshot{}, ErrInvalid
	}
	current.Status = approval.Expired
	current.Version++
	return current, nil
}

// EffectResult is internal to the eventual resource owner. APPROVED never
// means APPLIED; the executor must atomically record the actual effect.
func EffectResult(current Snapshot, version int64, result approval.Status, now time.Time) (Snapshot, error) {
	if err := validate(current, version, now); err != nil {
		return Snapshot{}, err
	}
	if current.Status != approval.Approved || current.CurrentStep != current.StepCount {
		return Snapshot{}, ErrInvalid
	}
	if !now.Before(current.ExpiresAt) {
		return Snapshot{}, ErrExpired
	}
	switch result {
	case approval.Applied, approval.Failed, approval.Stale:
		current.Status = result
	default:
		return Snapshot{}, ErrInvalid
	}
	current.Version++
	return current, nil
}
