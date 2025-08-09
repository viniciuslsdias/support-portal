package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/viniciuslsdias/support-portal/config"
	"github.com/viniciuslsdias/support-portal/internal/database"
	"github.com/viniciuslsdias/support-portal/internal/repository"
)

// Template cache
var templates *template.Template

func init() {
	// Parse all templates
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

// Home handler - displays the support request form
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := templates.ExecuteTemplate(w, "form.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		db := repository.New(database.GetPool())

		db.CreateTicket(context.Background(), repository.CreateTicketParams{
			FullName:            sql.NullString{r.FormValue("fullname"), true},
			EmailAddress:        sql.NullString{r.FormValue("email"), true},
			IssueCategory:       repository.Categories(strings.ToLower(r.FormValue("category"))),
			Priority:            repository.Priorities(strings.ToLower(r.FormValue("priority"))),
			IssueSummary:        sql.NullString{r.FormValue("summary"), true},
			DetailedDescription: sql.NullString{r.FormValue("description"), true},
			Department:          repository.Departments(r.FormValue("department")),
		})

		// Redirect to tickets page
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	}
}

// Tickets handler - displays all submitted tickets
func ticketsHandler(w http.ResponseWriter, r *http.Request) {

	db := repository.New(database.GetPool())

	data, _ := db.GetAllTickets(context.Background())

	err := templates.ExecuteTemplate(w, "tickets.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Ticket detail handler - displays individual ticket details
func ticketDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	db := repository.New(database.GetPool())

	data, err := db.GetTicket(context.Background(), id)

	if err != nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	err = templates.ExecuteTemplate(w, "detail.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	//Load data of config file
	cfg := config.GetConfig()

	// For gracefull shutdown
	ctx, _ := context.WithCancel(context.Background())

	// Initialize database connection
	database.ConnectPool(ctx, *cfg)
	defer database.ClosePool()

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tickets", ticketsHandler)
	http.HandleFunc("/ticket", ticketDetailHandler)

	log.Printf("Server starting on :%s", cfg.HTTPServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPServerPort), nil))
}
