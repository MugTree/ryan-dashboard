package dashboard

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"time"
// )

// const SessionTokenName string = "session_token"

// var sessions = map[string]Session{}

// type Session struct {
// 	Id       int
// 	UserName string
// 	Expiry   time.Time
// }

// func (s Session) IsExpired() bool {
// 	return s.Expiry.Before(time.Now())
// }

// type AuthKey struct{}
// type UserKey struct{}

// type AuthUser struct {
// 	ID       int
// 	UserName string
// 	LoggedIn bool
// }

// func (user AuthUser) Check() bool {
// 	return user.ID != 0 && user.LoggedIn
// }

// func authenticateUser(r *http.Request) (AuthUser, error) {

// 	c, err := r.Cookie(SessionTokenName)
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			return AuthUser{}, nil // most requests will follow this path
// 		}
// 		return AuthUser{}, nil
// 	}

// 	sessionToken := c.Value

// 	userSession, exists := sessions[sessionToken]

// 	if !exists {
// 		return AuthUser{}, nil
// 	}

// 	if userSession.IsExpired() {
// 		delete(sessions, sessionToken)
// 		return AuthUser{}, errors.New("session expired")
// 	}

// 	user := AuthUser{
// 		ID:       userSession.Id,
// 		UserName: userSession.UserName,
// 		LoggedIn: true,
// 	}

// 	return user, nil
// }

// func setUser(ctx context.Context, u *AuthUser) context.Context {
// 	return context.WithValue(ctx, UserKey{}, u)
// }

// func getUser(ctx context.Context) *AuthUser {
// 	user, ok := ctx.Value(UserKey{}).(*AuthUser)

// 	if !ok {
// 		return nil
// 	}

// 	return user
// }
