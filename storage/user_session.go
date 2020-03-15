// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage // import "miniflux.app/storage"

import (
	"database/sql"
	"fmt"

	"miniflux.app/crypto"
	"miniflux.app/model"
)

// UserSessions returns the list of sessions for the given user.
func (s *Storage) UserSessions(userID int64) (model.UserSessions, error) {
	query := `
		SELECT
			id,
			user_id,
			token,
			created_at,
			user_agent,
			ip
		FROM
			user_sessions
		WHERE
			user_id=$1 ORDER BY id DESC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch user sessions: %v`, err)
	}
	defer rows.Close()

	var sessions model.UserSessions
	for rows.Next() {
		var session model.UserSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Token,
			&session.CreatedAt,
			&session.UserAgent,
			&session.IP,
		)

		if err != nil {
			return nil, fmt.Errorf(`store: unable to fetch user session row: %v`, err)
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// CreateUserSession creates a new sessions.
func (s *Storage) CreateUserSession(username, userAgent, ip string) (sessionID string, userID int64, err error) {
	query := `SELECT id FROM users WHERE username = LOWER($1)`
	err = s.db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return "", 0, fmt.Errorf(`store: unable to fetch user ID: %v`, err)
	}

	token := crypto.GenerateRandomString(64)
	query = `INSERT INTO user_sessions (token, user_id, user_agent, ip) VALUES ($1, $2, $3, $4)`
	_, err = s.db.Exec(query, token, userID, userAgent, ip)
	if err != nil {
		return "", 0, fmt.Errorf(`store: unable to create user session: %v`, err)
	}

	return token, userID, nil
}

// UserSessionByToken finds a session by the token.
func (s *Storage) UserSessionByToken(token string) (*model.UserSession, error) {
	var session model.UserSession

	query := `
		SELECT
			id,
			user_id,
			token,
			created_at,
			user_agent,
			ip 
		FROM
			user_sessions
		WHERE
			token = $1
	`
	err := s.db.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.CreatedAt,
		&session.UserAgent,
		&session.IP,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch user session: %v`, err)
	default:
		return &session, nil
	}
}

// RemoveUserSessionByToken remove a session by using the token.
func (s *Storage) RemoveUserSessionByToken(userID int64, token string) error {
	query := `DELETE FROM user_sessions WHERE user_id=$1 AND token=$2`
	result, err := s.db.Exec(query, userID, token)
	if err != nil {
		return fmt.Errorf(`store: unable to remove this user session: %v`, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`store: unable to remove this user session: %v`, err)
	}

	if count != 1 {
		return fmt.Errorf(`store: nothing has been removed`)
	}

	return nil
}

// RemoveUserSessionByID remove a session by using the ID.
func (s *Storage) RemoveUserSessionByID(userID, sessionID int64) error {
	query := `DELETE FROM user_sessions WHERE user_id=$1 AND id=$2`
	result, err := s.db.Exec(query, userID, sessionID)
	if err != nil {
		return fmt.Errorf(`store: unable to remove this user session: %v`, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`store: unable to remove this user session: %v`, err)
	}

	if count != 1 {
		return fmt.Errorf(`store: nothing has been removed`)
	}

	return nil
}

// CleanOldUserSessions removes user sessions older than specified days.
func (s *Storage) CleanOldUserSessions(days int) int64 {
	query := `
		DELETE FROM
			user_sessions
		WHERE
			id IN (SELECT id FROM user_sessions WHERE created_at < now() - interval '%d days')
	`
	result, err := s.db.Exec(fmt.Sprintf(query, days))
	if err != nil {
		return 0
	}

	n, _ := result.RowsAffected()
	return n
}
