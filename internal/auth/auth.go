package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName = "session_id"
	SessionDuration   = 24 * time.Hour
)

// Session represents a user session
type Session struct {
	ID        string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore manages user sessions
type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*Session),
	}

	// Start cleanup goroutine
	go store.cleanup()

	return store
}

// Create creates a new session for a user
func (s *SessionStore) Create(username string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := generateSessionID()
	now := time.Now()

	session := &Session{
		ID:        sessionID,
		Username:  username,
		CreatedAt: now,
		ExpiresAt: now.Add(SessionDuration),
	}

	s.sessions[sessionID] = session
	return session
}

// Get retrieves a session by ID
func (s *SessionStore) Get(sessionID string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil
	}

	if time.Now().After(session.ExpiresAt) {
		return nil
	}

	return session
}

// Delete removes a session
func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// cleanup removes expired sessions periodically
func (s *SessionStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:])
}

// Authenticator handles authentication
type Authenticator struct {
	username     string
	passwordHash string
	sessions     *SessionStore
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(username, password string) *Authenticator {
	return &Authenticator{
		username:     username,
		passwordHash: hashPassword(password),
		sessions:     NewSessionStore(),
	}
}

// Validate checks if username and password are correct
func (a *Authenticator) Validate(username, password string) bool {
	return username == a.username && hashPassword(password) == a.passwordHash
}

// Login creates a session and sets cookie
func (a *Authenticator) Login(c *gin.Context, username, password string) bool {
	if !a.Validate(username, password) {
		return false
	}

	session := a.sessions.Create(username)
	c.SetCookie(
		SessionCookieName,
		session.ID,
		int(SessionDuration.Seconds()),
		"/",
		"",
		false, // secure - set to true in production with HTTPS
		true,  // httpOnly
	)

	return true
}

// Logout destroys the session
func (a *Authenticator) Logout(c *gin.Context) {
	sessionID, err := c.Cookie(SessionCookieName)
	if err == nil {
		a.sessions.Delete(sessionID)
	}

	c.SetCookie(SessionCookieName, "", -1, "/", "", false, true)
}

// GetSession returns the current session
func (a *Authenticator) GetSession(c *gin.Context) *Session {
	sessionID, err := c.Cookie(SessionCookieName)
	if err != nil {
		return nil
	}

	return a.sessions.Get(sessionID)
}

// AuthMiddleware is the Gin middleware for authentication
func (a *Authenticator) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := a.GetSession(c)
		if session == nil {
			// Check if it's an API request
			if isAPIRequest(c) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				c.Abort()
				return
			}

			// Redirect to login for web requests
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user", session.Username)
		c.Set("session", session)
		c.Next()
	}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func isAPIRequest(c *gin.Context) bool {
	// Check if path starts with /api
	return len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api"
}
