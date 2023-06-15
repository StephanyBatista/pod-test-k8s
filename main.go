package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	println("Starting server...")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	readinessProbeStatus := 200
	if os.Getenv("READINESS_PROBE_STATUS") != "" {
		readinessProbeStatus, _ = strconv.Atoi(os.Getenv("READINESS_PROBE_STATUS"))
	}
	r.Get("/readiness-probe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(readinessProbeStatus)
	})

	livenessProbeStatus := 200
	livenessProbeDuration := 0
	r.Get("/liveness-probe", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "" {
			//will change the status code of the liveness probe and so the container will be restarted
			livenessProbeStatus, _ = strconv.Atoi(r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("duration") != "" {
			//will change the duration of endpoint and so the container will be restarted
			livenessProbeDuration, _ = strconv.Atoi(r.URL.Query().Get("duration"))
		}
		time.Sleep(time.Duration(livenessProbeDuration) * time.Second)
		w.WriteHeader(livenessProbeStatus)
	})

	r.Get("/force-error", func(w http.ResponseWriter, r *http.Request) {
		//will force a error and the container will be killed
		_ = fib(1000000000000000000)
	})

	r.Get("/decode", func(w http.ResponseWriter, r *http.Request) {
		jsonString := r.URL.Query().Get("token")

		token, err := jwt.Parse(jsonString, nil)
		if err != nil {
			render.Status(r, 404)
			render.JSON(w, r, map[string]interface{}{"error": err.Error()})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		render.JSON(w, r, claims)
	})

	println("Server created...")

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3015"
	} else {
		port = ":" + port
	}
	println("Server started and listening on port " + port + "...")
	http.ListenAndServe(port, r)
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
