package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

type PageData struct {
	Content template.HTML
	Title   string
	Year    string
	IsMedia bool
	IsHomes bool
	Page    string
}

type ContactInfo struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Phone    string `form:"phone"`
	Message  string `form:"message"`
	Honeypot string `form:"honeypot"`
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/homes", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "views/homes.html", PageData{Page: "homes", IsHomes: true, Title: "Haier Homes | Kansas City Real Estate"})
	})

	http.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "views/media.html", PageData{Page: "media", IsMedia: true, Title: "Haier the Creator | Haier Media"})
	})

	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "views/privacy.html", PageData{Page: "homes", IsHomes: true, Title: "Privacy Policy | Haier Homes"})
	})

	http.HandleFunc("/terms", func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, "views/terms.html", PageData{Page: "homes", IsHomes: true, Title: "Terms & Conditions | Haier Homes"})
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Error parsing form")
			return
		}

		// Create an instance of the struct and fill it with form values
		data := ContactInfo{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Message:  r.FormValue("message"),
			Honeypot: r.FormValue("honeypot"),
		}

		if data.Honeypot != "" {
			fmt.Fprintf(w, "Error parsing form")
			return
		}

		mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), os.Getenv("MG_API_KEY"))
		sender := "Contact Form <noreply@haiertherealtor.com>"
		subject := "New Contact!"
		emailBody := fmt.Sprintf(
			"New Message from %s\n%s\n%s\n\n%s",
			data.Name,
			data.Email,
			data.Phone,
			data.Message,
		)
		recipient := os.Getenv("RECIPIENT_EMAIL")

		// The message object allows you to add attachments and Bcc recipients
		message := mg.NewMessage(sender, subject, emailBody, recipient)

		mgCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// Send the message with a 10 second timeout
		resp, id, err := mg.Send(mgCtx, message)

		if err != nil {
			fmt.Fprintf(w, "Error sending email")
			return
		}

		fmt.Printf("ID: %s Resp: %s\n", id, resp)

		fmt.Fprintf(w, "Thank you for reaching out! I will get back to you soon.")
		return
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

func renderPage(w http.ResponseWriter, tmpl string, data PageData) {
	renderTemplate(w, "views/layout.html", PageData{
		IsHomes: data.IsHomes,
		IsMedia: data.IsMedia,
		Page:    data.Page,
		Title:   data.Title,
		Content: renderTemplateToString(tmpl, data),
		Year:    fmt.Sprint(time.Now().Year()),
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateToString(tmpl string, data PageData) template.HTML {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return ""
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)
	if err != nil {
		return ""
	}

	return template.HTML(buffer.String())
}
