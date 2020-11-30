package crud

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ssd-coursework/routes/user"
	"strconv"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type M map[string]interface{}

type Employee struct {
	Id   int
	Name string
	City string
}

type Ticket struct {
	TicketID     int
	Title        string
	TicketType   string
	Status       string
	CreatedAt    string
	TicketNumber int
	Description  string
	FinderID     string
	AssignedTo   string
	Priority     string
}

type Comment struct {
	CommentID   int
	TicketID    int
	UserID      string
	TimeStamp   string
	Description string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "admin"
	dbPass := "admin"
	dbName := "bug_tracker"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))
var currentTicketID int

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM TICKET ORDER BY ticketID DESC")
	if err != nil {
		panic(err.Error())
	}
	ticket := Ticket{}
	res := []Ticket{}
	for selDB.Next() {
		var ticketID, ticketNumber int
		var title, ticketType, status, createdAt, description, finderID, assignedTo, priority string
		err = selDB.Scan(&ticketID, &title, &ticketType, &status, &createdAt, &ticketNumber, &description, &finderID, &assignedTo, &priority)
		if err != nil {
			panic(err.Error())
		}
		ticket.TicketID = ticketID
		ticket.Title = title
		ticket.TicketType = ticketType
		ticket.Status = status
		ticket.CreatedAt = createdAt
		ticket.TicketNumber = ticketNumber
		ticket.Description = description
		ticket.FinderID = user.GetUsersNameFromAuth(w, r, finderID)
		ticket.AssignedTo = user.GetUsersNameFromAuth(w, r, assignedTo)
		ticket.Priority = priority
		res = append(res, ticket)
	}
	err = tmpl.ExecuteTemplate(w, "Index", res)
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
}

// func Index(w http.ResponseWriter, r *http.Request) {
// 	db := dbConn()
// 	selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	emp := Employee{}
// 	res := []Employee{}
// 	for selDB.Next() {
// 		var id int
// 		var name, city string
// 		err = selDB.Scan(&id, &name, &city)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		emp.Id = id
// 		emp.Name = name
// 		emp.City = city
// 		res = append(res, emp)
// 	}
// 	err = tmpl.ExecuteTemplate(w, "Index", res)
// 	if err != nil {
// 		log.Print(err.Error())
// 	}
// 	defer db.Close()
// }

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	// ToDo https://johnweldon.com/blog/quick-tip-remove-query-param-from-url-in-go/
	nID := r.URL.Query().Get("ticketID")
	selDB, err := db.Query("SELECT * FROM TICKET WHERE ticketID=?", nID)
	if err != nil {
		panic(err.Error())
	}
	ticket := Ticket{}
	for selDB.Next() {
		var ticketID, ticketNumber int
		var title, ticketType, status, createdAt, description, finderID, assignedTo, priority string
		err = selDB.Scan(&ticketID, &title, &ticketType, &status, &createdAt, &ticketNumber, &description, &finderID, &assignedTo, &priority)
		if err != nil {
			panic(err.Error())
		}
		ticket.TicketID = ticketID
		ticket.Title = title
		ticket.TicketType = ticketType
		ticket.Status = status
		ticket.CreatedAt = createdAt
		ticket.TicketNumber = ticketNumber
		ticket.Description = description
		ticket.FinderID = user.GetUsersNameFromAuth(w, r, finderID)
		ticket.AssignedTo = user.GetUsersNameFromAuth(w, r, assignedTo)
		ticket.Priority = priority
		currentTicketID = ticketID
	}
	fmt.Println(ticket)
	comments := getComments(w, r, ticket.TicketID)
	err = tmpl.ExecuteTemplate(w, "Show", M{
		"ticket":   ticket,
		"comments": comments,
	})
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
}

func getComments(w http.ResponseWriter, r *http.Request, ticketID int) []Comment {
	db := dbConn()
	// ToDo https://johnweldon.com/blog/quick-tip-remove-query-param-from-url-in-go/
	nID := ticketID
	selDB, err := db.Query("SELECT * FROM COMMENT WHERE ticketID=?", nID)
	if err != nil {
		panic(err.Error())
	}
	comment := Comment{}
	res := []Comment{}
	for selDB.Next() {
		var commentID, ticketID int
		var userID, description, timeStamp string
		err = selDB.Scan(&commentID, &ticketID, &userID, &timeStamp, &description)
		if err != nil {
			panic(err.Error())
		}
		comment.CommentID = commentID
		comment.TicketID = ticketID
		comment.UserID = user.GetUsersNameFromAuth(w, r, userID)
		comment.TimeStamp = timeStamp
		comment.Description = description
		res = append(res, comment)
	}

	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	return res
}

// AddComment Adds a comment to a ticket
func AddComment(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(user.GetSessionUsername(w, r))
	db := dbConn()
	if r.Method == "POST" {
		userID := user.GetSessionUsername(w, r)
		ticketID := currentTicketID
		description := r.FormValue("description")
		insForm, err := db.Prepare("INSERT INTO Comment(userID, ticketID, comment) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = insForm.Exec(userID, ticketID, description)
		if err != nil {
			log.Print(err.Error())
		}
		log.Println("INSERT: UserID: " + userID + " | ticketID: " + strconv.Itoa(ticketID) + " | description: " + description)
	}
	defer db.Close()
	urlRedirect := "/show?ticketID=" + strconv.Itoa(currentTicketID)
	http.Redirect(w, r, urlRedirect, 301)
}

func New(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "New", nil)
	if err != nil {
		log.Print(err.Error())
	}
}

// Edit is called when the user clicks on the edit button
// will load the values of the ticket ID into the form
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("ticketID")
	selDB, err := db.Query("SELECT * FROM TICKET WHERE ticketID=?", nId)
	if err != nil {
		panic(err.Error())
	}
	ticket := Ticket{}
	for selDB.Next() {
		var ticketID, ticketNumber int
		var title, ticketType, status, createdAt, description, finderID, assignedTo, priority string
		err = selDB.Scan(&ticketID, &title, &ticketType, &status, &createdAt, &ticketNumber, &description, &finderID, &assignedTo, &priority)
		if err != nil {
			panic(err.Error())
		}
		ticket.TicketID = ticketID
		ticket.Title = title
		ticket.TicketType = ticketType
		ticket.Status = status
		ticket.CreatedAt = createdAt
		ticket.TicketNumber = ticketNumber
		ticket.Description = description
		ticket.FinderID = user.GetUsersNameFromAuth(w, r, finderID)
		ticket.AssignedTo = user.GetUsersNameFromAuth(w, r, assignedTo)
		ticket.Priority = priority
		// currentTicketID = ticketID
	}
	fmt.Println(ticket)
	err = tmpl.ExecuteTemplate(w, "Edit", ticket)
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
}

// UpdateTicket is called when the user clicks on the submit button of the update form
func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		title := r.FormValue("title")
		ticketType := r.FormValue("ticketType")
		status := r.FormValue("status")
		ticketNumber := r.FormValue("ticketNumber")
		description := r.FormValue("description")
		finderID := user.GetUserIDFromName(w, r, r.FormValue("finderID"))
		assignedTo := user.GetUserIDFromName(w, r, r.FormValue("assignedTo"))
		priority := r.FormValue("priority")
		ticketID := r.FormValue("uid")

		insForm, err := db.Prepare(`UPDATE ticket SET Title=?, type=?, status=?, ticketNumber=?, description=?,
				finderID=?, assignedTo=?, priority=? WHERE ticketID=?`)
		if err != nil {
			panic(err.Error())
		}
		_, err = insForm.Exec(title, ticketType, status, ticketNumber, description, finderID, assignedTo, priority, ticketID)
		if err != nil {
			log.Print(err.Error())
		}
		log.Println("UPDATE: TICKET ID: " + ticketID + " | Ticket: " + title + " | ticketType: " + ticketType)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Insert UPDATE THIS TO CREATE NEW TICKET
func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = insForm.Exec(name, city)
		if err != nil {
			log.Print(err.Error())
		}
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

// Delete REMOVE DELETE, DON'T THINK WE NEED THIS ANYMORE
func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	_, err = delForm.Exec(emp)
	if err != nil {
		log.Print(err.Error())
	}
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
