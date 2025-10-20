package routes

import (
    "authentication/services"
    "crypto/rand"
    "encoding/base64"
    "net/http"
    "time"
)

// RegisterHTTPRoutes sets up net/http routes on the provided mux.
func RegisterHTTPRoutes(mux *http.ServeMux) {
    // Health/ping endpoint
    mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(`{"message":"pong"}`))
    })

    // Google OAuth: start login
    mux.HandleFunc("/api/oauth/google/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Generate state token and store in a short-lived cookie
        state := generateState()
        http.SetCookie(w, &http.Cookie{
            Name:     "oauth_state",
            Value:    state,
            Path:     "/",
            HttpOnly: true,
            Secure:   isHTTPS(r),
            SameSite: http.SameSiteLaxMode,
            Expires:  time.Now().Add(5 * time.Minute),
        })

        url := services.GetGoogleLoginURL(state)
        http.Redirect(w, r, url, http.StatusFound)
    })

    // Google OAuth: callback
    mux.HandleFunc("/api/oauth/google/callback", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }

        // Validate state
        stateCookie, err := r.Cookie("oauth_state")
        if err != nil || stateCookie.Value == "" {
            http.Error(w, "missing state", http.StatusBadRequest)
            return
        }
        returnedState := r.URL.Query().Get("state")
        if returnedState == "" || returnedState != stateCookie.Value {
            http.Error(w, "invalid state", http.StatusBadRequest)
            return
        }

        code := r.URL.Query().Get("code")
        if code == "" {
            http.Error(w, "missing code", http.StatusBadRequest)
            return
        }

        user, err := services.HandleGoogleCallback(code)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // For simplicity, return minimal JSON response
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("{\"id\":\"" + user.ID.String() + "\",\"email\":\"" + user.Email + "\"}"))
    })
}

func generateState() string {
    buf := make([]byte, 16)
    _, _ = rand.Read(buf)
    return base64.RawURLEncoding.EncodeToString(buf)
}

func isHTTPS(r *http.Request) bool {
    if r.TLS != nil {
        return true
    }
    // Support reverse proxy via X-Forwarded-Proto
    if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
        return true
    }
    return false
}
