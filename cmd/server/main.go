package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Ticket represents a support ticket
type Ticket struct {
	ID          int
	FullName    string
	Email       string
	Category    string
	Priority    string
	Summary     string
	Description string
	CreatedAt   time.Time
}

// In-memory storage for tickets
var tickets []Ticket
var ticketCounter int

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

		// Create new ticket
		ticketCounter++
		ticket := Ticket{
			ID:          ticketCounter,
			FullName:    r.FormValue("fullname"),
			Email:       r.FormValue("email"),
			Category:    r.FormValue("category"),
			Priority:    r.FormValue("priority"),
			Summary:     r.FormValue("summary"),
			Description: r.FormValue("description"),
			CreatedAt:   time.Now(),
		}

		// Add to tickets slice
		tickets = append(tickets, ticket)

		// Redirect to tickets page
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	}
}

// Tickets handler - displays all submitted tickets
func ticketsHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Tickets []Ticket
	}{
		Tickets: tickets,
	}

	err := templates.ExecuteTemplate(w, "tickets.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Ticket detail handler - displays individual ticket details
func ticketDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	var ticket *Ticket
	for i := range tickets {
		if tickets[i].ID == id {
			ticket = &tickets[i]
			break
		}
	}

	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	err = templates.ExecuteTemplate(w, "detail.html", ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tickets", ticketsHandler)
	http.HandleFunc("/ticket", ticketDetailHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
