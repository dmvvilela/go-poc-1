package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models" // models package where Contact schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to convert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// Read .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// return the connection
	return db
}

// CreateContact create a contact in the postgres db
func CreateContact(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty contact of type models.Contact
	var contact models.Contact

	// decode the json request to contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert contact function and pass the contact
	insertID := insertContact(contact)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Contact created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetContact will return a single contact by its id
func GetContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get the contactid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the getContact function with contact id to retrieve a single contact
	contact, err := getContact(int64(id))
	if err != nil {
		log.Fatalf("Unable to get contact. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(contact)
}

// GetAllContacts will return all the contacts
func GetAllContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all the contacts in the db
	contacts, err := getAllContacts()
	if err != nil {
		log.Fatalf("Unable to get all contact. %v", err)
	}

	// send all the contacts as response
	json.NewEncoder(w).Encode(contacts)
}

// UpdateContact update contact's detail in the postgres db
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the contactid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// create an empty contact of type models.Contact
	var contact models.Contact

	// decode the json request to contact
	err = json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update contact to update the contact
	updatedRows := updateContact(int64(id), contact)

	// format the message string
	msg := fmt.Sprintf("Contact updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteContact delete contact's detail in the postgres db
func DeleteContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the contactid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteContact, convert the int to int64
	deletedRows := deleteContact(int64(id))

	// format the message string
	msg := fmt.Sprintf("Contact updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions -------------------------
// insert one contact in the DB
func insertContact(contact models.Contact) int64 {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning contactid will return the id of the inserted contact
	sqlStatement := `INSERT INTO contacts (name, email) VALUES ($1, $2) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, contact.Name, contact.Email).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one contact from the DB by its contactid
func getContact(id int64) (models.Contact, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a contact of models.Contact type
	var contact models.Contact

	// create the select sql query
	sqlStatement := `SELECT * FROM contacts WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to contact
	err := row.Scan(&contact.ID, &contact.Name, &contact.Email)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return contact, nil
	case nil:
		return contact, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty contact on error
	return contact, err
}

// get one contact from the DB by its contactid
func getAllContacts() ([]models.Contact, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var contacts []models.Contact

	// create the select sql query
	sqlStatement := `SELECT * FROM contacts`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var contact models.Contact

		// unmarshal the row object to contact
		err = rows.Scan(&contact.ID, &contact.Name, &contact.Email)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the contact in the contacts slice
		contacts = append(contacts, contact)

	}

	// return empty contact on error
	return contacts, err
}

// update contact in the DB
func updateContact(id int64, contact models.Contact) int64 {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE contacts SET name=$2, email=$3 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, contact.Name, contact.Email)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete contact in the DB
func deleteContact(id int64) int64 {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM contacts WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
