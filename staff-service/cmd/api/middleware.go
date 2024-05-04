package main

import (
	"log"
	"net/http"
	"staff-service/internal/data"
	"strings"
)

func (app *Config) requireActivatedUser(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.contextGetUser(r)
		if err != nil {
			app.errorJson(w, err, http.StatusUnauthorized)
			return
		}
		if !user.Activated {
			app.errorJson(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Config) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		log.Println("Authorization header: ", authorizationHeader)
		if authorizationHeader == "" {
			app.errorJson(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.errorJson(w, ErrAuthentication, http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		log.Println("Token: ", token)

		user, err := app.AuthenticateViaGrpc(token)
		if err != nil {
			app.rpcErrorJson(w, err)
			return
		}

		r = app.contextSetUserId(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *Config) requireOwnerRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.contextGetUser(r)
		if err != nil {
			app.errorJson(w, err, http.StatusUnauthorized)
			return
		}
		if !(user.Type == "owner") {
			app.errorJson(w, ErrAuthentication, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Config) requireInstOwnerForOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.contextGetUser(r)
		if err != nil {
			app.errorJson(w, err, http.StatusUnauthorized)
			return
		}
		var input data.Employee
		err = app.readJSON(w, r, &input)
		if err != nil {
			app.errorJson(w, err, http.StatusBadRequest)
			return
		}

		institution, err := app.GetInstitution(input.InstId)
		if err != nil {
			app.rpcErrorJson(w, err)
			return
		}
		if institution.OwnerId != user.ID {
			app.errorJson(w, ErrBadRequest, http.StatusForbidden)
			return
		}
		r = app.contextSetEmployee(r, &input)
		next.ServeHTTP(w, r)
	})
}

// func (app *Config) rateLimit(next http.Handler) http.Handler {

// 	type client struct {
// 		limiter  *rate.Limiter
// 		lastSeen time.Time
// 	}
// 	var (
// 		mu      sync.Mutex
// 		clients = make(map[string]*client)
// 	)
// 	go func() {
// 		for {
// 			time.Sleep(time.Minute)
// 			mu.Lock()

// 			for ip, client := range clients {
// 				if time.Since(client.lastSeen) > 3*time.Minute {
// 					delete(clients, ip)
// 				}
// 			}
// 			mu.Unlock()
// 		}
// 	}()
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if app.config.limiter.enabled {
// 			ip, _, err := net.SplitHostPort(r.RemoteAddr)
// 			if err != nil {
// 				app.serverErrorResponse(w, r, err)
// 				return
// 			}
// 			mu.Lock()
// 			if _, found := clients[ip]; !found {
// 				clients[ip] = &client{
// 					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
// 				}
// 			}

// 			clients[ip].lastSeen = time.Now()
// 			if !clients[ip].limiter.Allow() {
// 				mu.Unlock()
// 				app.rateLimitExceededResponse(w, r)
// 				return
// 			}
// 			mu.Unlock()
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
