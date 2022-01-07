package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/philipsahli/unfollow/pkg/twitter"
	"github.com/xeonx/timeago"
)

var homepageTpl *template.Template

type TemplateData struct {
	Users []twitter.User
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var err error
	homepageTpl, err = template.New("home.gohtml").Funcs(template.FuncMap{
		"timeago": func(t time.Time) string {
			return timeago.English.Format(t)
		},
	}).ParseFiles("pkg/web/home.gohtml")
	if err != nil {
		panic(err)
	}

	users := twitter.GetUsers()

	data := TemplateData{
		Users: users,
	}

	homepageTpl.Execute(w, data)
}

func Start() {
	var wait time.Duration

	r := mux.NewRouter()
	// Add your routes as needed

	r.HandleFunc("/", HomeHandler)

	srv := &http.Server{
		Addr: "0.0.0.0:8081",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
