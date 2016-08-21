package ggdb

import (
	"database/sql"
	"log"
	"os"
	"os/user"

	// Import the SQLite driver.
	_ "github.com/mattn/go-sqlite3"

	"github.com/leerumler/gengar/ggconf"
)

// Phrase holds
type Phrase struct {
	Phrase    string
	Expansion Expansion
}

// Expansion holds
type Expansion struct {
	Name, Expansion string
	ID              int
}

// Category holds
type Category struct {
	ID   int
	Name string
}

// GGDB holds a connection to gengar's database.
var GGDB *sql.DB

// findGGDB locates the gengar database file.
func findGGDB() *string {

	// Check the current user.
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Find the homedir and add the rest of the path.
	homedir := usr.HomeDir

	// Make the gengar config directory if it doesn't exist.
	confdir := homedir + "/.config/gengar"
	os.MkdirAll(confdir, 0700)

	// Append database name and return a pointer.
	dbfile := confdir + "/gg.db"
	return &dbfile
}

// ConnectGGDB establishes a connection to Gengar's database.
func connectGGDB() *sql.DB {

	// Find the database and open it.
	dbfile := findGGDB()
	db, err := sql.Open("sqlite3", *dbfile)
	if err != nil {
		log.Fatal(err)
	}

	// Return a pointer to that database connection.
	return db
}

// CleanSlate creates an empty SQLite database to act as gengar's central data store.
func CleanSlate() {

	// Get a pointer to the database connection.
	db := connectGGDB()
	defer db.Close()

	// Drop existing tables and (re-)create them.
	createTables := `
	DROP TABLE IF EXISTS categories;
	DROP TABLE IF EXISTS expansions;
	DROP TABLE IF EXISTS phrases;
	CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	INSERT INTO categories (name) VALUES ("default");
	CREATE TABLE expansions (
		cat_id INTEGER DEFAULT 1,
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		expansion TEXT NOT NULL UNIQUE,
		FOREIGN KEY (cat_id) REFERENCES categories(id)
		);
	CREATE TABLE phrases (
		exp_id INTEGER,
		name TEXT PRIMARY KEY,
		FOREIGN KEY (exp_id) REFERENCES expansions(id)
	);

	`
	_, err := db.Exec(createTables)

	// Die on error.
	if err != nil {
		log.Fatal("Couldn't create CleanSlate:", err)
	}
}

// AddCategory inserts a new Category into gengar's database.
func AddCategory(cat *Category) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the category.
	_, err := db.Exec("INSERT INTO categories (name) VALUES ($1);", cat.Name)
	if err != nil {
		log.Fatal("Couldn't insert category:", err)
	}
}

// AddExpansion inserts a new Expansion into gengar's database.
func AddExpansion(exp *Expansion) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the expansion.
	_, err := db.Exec("INSERT INTO expansions (name, expansion) VALUES ($1, $2);", exp.Name, exp.Expansion)
	if err != nil {
		log.Fatal("Couldn't Insert Expansions:", err)
	}
}

// AddPhrase maps a phrase to an expansion in gengar's database.
func AddPhrase(phrase *Phrase) {

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Insert the phrase.
	_, err := db.Exec("INSERT INTO phrases (name, exp_id) VALUES ($1, $2);", phrase.Phrase, phrase.Expansion.ID)
	if err != nil {
		log.Fatal("Couldn't insert phrases:", err)
	}
}

// ReadCategories reads the available categories from the database.
func ReadCategories() []Category {

	// Create an empty slice of Categories to fill.
	var cats []Category

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for all available categories.
	rows, err := db.Query("SELECT id, name FROM categories;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Load the query's results in to new Categories and append them to the slice.
	for rows.Next() {
		var cat Category
		err = rows.Scan(&cat.ID, &cat.Name)
		cats = append(cats, cat)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Return the filled slice of Categories
	return cats
}

// ReadExpansions finds all of the Expansions within a given category.
func ReadExpansions(cat Category) []Expansion {

	// Create an empty slice of Expansions to fill.
	var exps []Expansion

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for Expansions matching the category's ID.
	rows, err := db.Query("SELECT id, name, expansion FROM expansions WHERE cat_id=$1;", cat.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Load the query's results in to new Expansions and append those to the slice.
	for rows.Next() {
		var exp Expansion
		err = rows.Scan(&exp.ID, &exp.Name, &exp.Expansion)
		exps = append(exps, exp)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Return the populated slice of Expansions.
	return exps
}

// ReadPhrases finds all of the Phrases mapped to a given Expansion.
func ReadPhrases(exp Expansion) []Phrase {

	// Create an empty slice of Phrases.
	var phrases []Phrase

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for Phrases matching the expansion's ID.
	rows, err := db.Query("SELECT name FROM phrases WHERE exp_ID=$1;", exp.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Load the query's results in to the Phrase slice.
	for rows.Next() {
		var phrase Phrase
		err = rows.Scan(&phrase.Phrase)
		phrase.Expansion = exp
		phrases = append(phrases, phrase)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Return the populated slice of phrases.
	return phrases
}

// ReadAllExpansions reads all of the available expansions from the database.
func ReadAllExpansions() []Expansion {

	// Create an empty slice of Expansions to fill.
	var exps []Expansion

	// Connect to the database.
	db := connectGGDB()
	defer db.Close()

	// Query the database for all Expansions.
	rows, err := db.Query("SELECT id, name, expansion FROM expansions;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Load the query's results in to new Expansions and append those to the slice.
	for rows.Next() {
		var exp Expansion
		err = rows.Scan(&exp.ID, &exp.Name, &exp.Expansion)
		exps = append(exps, exp)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return exps
}

// CreateTestDB creates a test database.
func CreateTestDB() {

	// Create some testing expansions.
	var exps []Expansion
	exps = append(exps, Expansion{Name: "Test 1", Expansion: "This is test 1!", ID: 1})
	exps = append(exps, Expansion{Name: "Test 2", Expansion: "this is test 2?", ID: 2})
	exps = append(exps, Expansion{Name: "Test 3", Expansion: "this is test 3!?@$", ID: 3})

	var phrases []Phrase
	phrases = append(phrases, Phrase{Phrase: "test1", Expansion: exps[0]})
	phrases = append(phrases, Phrase{Phrase: "test2", Expansion: exps[1]})
	phrases = append(phrases, Phrase{Phrase: "test3", Expansion: exps[2]})

	// Wipe/create a blank gengar database.
	CleanSlate()

	// Insert each of our testing expansions.
	for _, exp := range exps {
		AddExpansion(&exp)
	}
	for _, phrase := range phrases {
		AddPhrase(&phrase)
	}
}

// findExpansionID finds and returns the expansion ID from gengar's database.
// func findExpansionID(exp *Expansion) int {
//
// 	// Connect to the database.
// 	db := connectGGDB()
// 	defer db.Close()
//
// 	// Check the expansion ID in the database.
// 	var expID int
// 	err := db.QueryRow("SELECT id FROM expansions WHERE name=$1;", exp.Name).Scan(&expID)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		log.Fatal("No matching expansions found.")
// 	case err != nil:
// 		log.Fatal(err)
// 	}
//
// 	// Return the ID.
// 	return expID
// }

// ReadExpanders reads expansions from the database and returns a slice of ggconf.Expanders.
func ReadExpanders() []ggconf.Expander {
	var exps []ggconf.Expander

	// Get pointer to database connection.
	db := connectGGDB()
	defer db.Close()

	// Query the database for the expansions.
	rows, err := db.Query("SELECT exp_id, phrases.name, expansion FROM phrases JOIN expansions ON phrases.exp_id = expansions.id;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Read through the results and populate exps with returned info.
	for rows.Next() {
		var exp ggconf.Expander
		err = rows.Scan(&exp.ID, &exp.Phrase, &exp.Expansion)
		exps = append(exps, exp)
	}

	// Die on error.
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Return pointer to exps.
	return exps
}
