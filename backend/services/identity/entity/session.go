package entity

import "time"

type SessionId string
type SessionStatus string

const (
	SessionStatusActive  SessionStatus = "ACTIVE"
	SessionStatusRevoked SessionStatus = "REVOKED"
	SessionStatusExpired SessionStatus = "EXPIRED"
)

type Session struct {
	sessionId        SessionId
	accountId        AccountId
	refreshTokenHash string
	status           SessionStatus
	expiresAt        time.Time
	lastSeenAt       time.Time
	revokedAt        *time.Time
	createdAt        time.Time
	updatedAt        time.Time
}

func NewSession(
	sessionId SessionId,
	accountId AccountId,
	refreshTokenHash string,
	status SessionStatus,
	expiresAt time.Time,
	lastSeenAt time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Session {
	return Session{
		sessionId:        sessionId,
		accountId:        accountId,
		refreshTokenHash: refreshTokenHash,
		status:           status,
		expiresAt:        expiresAt,
		lastSeenAt:       lastSeenAt,
		revokedAt:        copyTime(revokedAt),
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

func (session Session) GetSessionId() SessionId {
	return session.sessionId
}

func (session Session) GetAccountId() AccountId {
	return session.accountId
}

func (session Session) GetRefreshTokenHash() string {
	return session.refreshTokenHash
}

func (session Session) GetStatus() SessionStatus {
	return session.status
}

func (session Session) GetExpiresAt() time.Time {
	return session.expiresAt
}

func (session Session) GetLastSeenAt() time.Time {
	return session.lastSeenAt
}

func (session Session) GetRevokedAt() *time.Time {
	return copyTime(session.revokedAt)
}

func (session Session) GetCreatedAt() time.Time {
	return session.createdAt
}

func (session Session) GetUpdatedAt() time.Time {
	return session.updatedAt
}

func (session Session) SetRefreshTokenHash(refreshTokenHash string) Session {
	session.refreshTokenHash = refreshTokenHash
	return session
}

func (session Session) SetStatus(status SessionStatus) Session {
	session.status = status
	return session
}

func (session Session) SetLastSeenAt(lastSeenAt time.Time) Session {
	session.lastSeenAt = lastSeenAt
	return session
}

func (session Session) SetRevokedAt(revokedAt *time.Time) Session {
	session.revokedAt = copyTime(revokedAt)
	return session
}

func (session Session) SetUpdatedAt(updatedAt time.Time) Session {
	session.updatedAt = updatedAt
	return session
}
