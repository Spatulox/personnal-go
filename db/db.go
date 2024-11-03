package db

import (
	"database/sql"
	"errors"
	"fmt"
	. "github.com/Spatulox/personnal-go/log"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
)

type Db struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDb(host string, port int, user, password, dbName string) *Db {
	return &Db{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}
}

// Fonction qui créer la connexion avec la BDD Mysql
// Renvoie une structure database
func connectDB(d *Db) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		Log.Error("Something went wrong when connecting to the database", err)
		return nil, err
	}

	// Vérifier la connexion
	err = db.Ping()
	if err != nil {
		Log.Error("Impossible to ping the database", err)
		return nil, err
	}

	Log.Infos("Connexion successful")
	return db, nil
}

//
// ------------------------------------------------------------------------------------------------ //
//

// SelectDB Fonction pour sélectionner des données dans une table SQL
// Paramètres :
// - Table SQL - String |
// - Colonnes où rentrer les données - Tableau de string |
// - Possiblité de rajouter une chaine de caractère pour les selects avec des INNER JOIN / RIGHT JOIN / ETC... - String |
// - Condition si besoin - String |
// - Voir la requête - Boolean |
func (d *Db) SelectDB(table string, column []string, join *string, condition *string, debug ...bool) ([]map[string]interface{}, error) {

	var err error = errors.New("Some error occurred")

	if checkData(table, column, nil, condition) == false {
		Log.Error("Plz check your parameters")
		return nil, err
	}

	if join != nil && reflect.TypeOf(*join).Kind() != reflect.String {
		Log.Error("Innerjoin value need to be a string")
		return nil, err
	}

	var db, errC = connectDB(d)

	if errC != nil {
		return nil, err
	}

	defer db.Close()

	if db == nil {
		Log.Error("What da heck bro, db is null ??")
		return nil, err
	}

	// checking the right format
	var columns = ArrayToString(column, true)

	if columns == "" {
		Log.Error("Impossible to transform the columns array into a string")
		return nil, err
	}

	var query *sql.Rows
	var queryString string

	// Créer la chaîne de caractère pour la requête, en fonction des paramètres passés
	if condition == nil && join == nil {
		query, err = db.Query("SELECT " + columns + " FROM " + table)
		queryString = "SELECT " + columns + " FROM " + table
		if err != nil {
			Log.Error("Erreur : ", err)
			Log.Debug(queryString)
			return nil, err
		}
	} else if condition != nil && join == nil {
		query, err = db.Query("SELECT " + columns + " FROM " + table + " WHERE " + *condition)
		queryString = "SELECT " + columns + " FROM " + table + " WHERE " + *condition
		if err != nil {
			Log.Error("Error : ", err)
			Log.Debug(queryString)
			return nil, err
		}
	} else {
		query, err = db.Query("SELECT " + columns + " FROM " + table + " " + *join + " WHERE " + *condition)
		queryString = "SELECT " + columns + " FROM " + table + " " + *join + " WHERE " + *condition
		if err != nil {
			Log.Error("Error : ", err)
			Log.Debug(queryString)
			return nil, err
		}
	}

	// si le paramètre debug existe
	if len(debug) > 0 && debug[0] {
		Log.Debug(queryString)
	}

	// Appelle un fonction qui execute la query et la transforme en "objet json"
	var result = transformQueryToMap(query)

	if err := query.Err(); err != nil {
		Log.Error("An error Occured : ", err)
		return nil, err
	}

	return result, nil
}

//
// ------------------------------------------------------------------------------------------------ //
//

// InsertDB Fonction pour insérer des données dans une table SQL
// Paramètres :
// - Table SQL - String |
// - Colonnes où rentrer les données - Tableau de string |
// - Valeur des données - Tableau de string |
// - Voir la requête - Boolean |
func (d *Db) InsertDB(table string, column []string, value []string, debug ...bool) {

	if checkData(table, column, value, nil) == false {
		return
	}

	var db, errC = connectDB(d)

	if errC != nil {
		return
	}

	defer db.Close()

	if db == nil {
		Log.Error("What da heck bro, db is null ??")
		return
	}

	// Appelle d'une fonction pour transformer les colonnes en chaine de caractère unique
	var columns = ArrayToString(column, true)

	// Appelle d'une fonction pour transformer les valeurs en chaine de caractère unique
	// Rajoute des quotes si besoin pour les chaines de caractères, rien si c'est un nombre
	var values = ArrayToString(value)

	if columns == "" {
		Log.Error("Impossible to transform the columns array into a string")
		return
	}

	if values == "" {
		Log.Error("Impossible to transform the columns array into a string")
		return
	}

	var query *sql.Rows
	var queryString string
	var err error

	// Execute la requête
	query, err = db.Query("INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ")")
	queryString = "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ")"
	if err != nil {
		Log.Error("Error : ", err)
		Log.Debug(queryString)
		return
	}

	if err := query.Err(); err != nil {
		Log.Error("An error Occured : ", err)
		return
	}

	if len(debug) > 0 && debug[0] {
		Log.Debug(queryString)
	}

	return

}

//
// ------------------------------------------------------------------------------------------------ //
//

// UpdateDB Fonction pour mettre à jour une table SQL
// Paramètres :
// - Table SQL - String |
// - Colonnes où rentrer les données - Tableau de string |
// - Valeur des données à rentrer - Tableau de string |
// - Condition si besoin - String |
// - Voir la requête - Boolean |
func (d *Db) UpdateDB(table string, column []string, value []string, condition *string, debug ...bool) {

	if checkData(table, column, value, condition) == false {
		return
	}

	if condition == nil {
		Log.Error("Plz enter a condition to update the table. If you don't want to enter condition put a \"-1\" instead")
		return
	}

	var db, errC = connectDB(d)

	if errC != nil {
		return
	}

	defer db.Close()

	if db == nil {
		Log.Error("What da heck bro, db is null ??")
		return
	}

	var query *sql.Rows
	var queryString string
	var err error

	var set = ConcatColumnWithValues(column, value)

	if set == "" {
		return
	}

	if condition != nil {
		query, err = db.Query("UPDATE " + table + " SET " + set + " WHERE " + *condition)
		queryString = "UPDATE " + table + " SET " + set + " WHERE " + *condition
		if err != nil {
			Log.Error("Error : ", err)
			Log.Debug(queryString)
			return
		}
	} else if *condition == "-1" {
		query, err = db.Query("UPDATE " + table + " SET " + set)
		queryString = "UPDATE " + table + " SET " + set
		if err != nil {
			Log.Error("Error : ", err)
			Log.Debug(queryString)
			return
		}
	}

	if err := query.Err(); err != nil {
		Log.Error("An error Occured : ", err)
		return
	}

	if len(debug) > 0 && debug[0] {
		Log.Debug(queryString)
	}

	return

}

//
// ------------------------------------------------------------------------------------------------ //
//

// DeleteDB Fonction pour supprimer les données d'une table SQL
// Paramètres :
// - Condition si besoin - String |
// - Voir la requête - Boolean |
func (d *Db) DeleteDB(table string, condition *string, debug ...bool) {
	// DELETE FROM table WHERE condition

	if reflect.TypeOf(table) != reflect.TypeOf("") || table == "" {
		Log.Error("Need a table name")
	}

	if condition == nil {
		Log.Error("Plz enter a condition to delete a row from a the table. If you don't want to enter condition put a \"-1\" instead")
		return
	}

	var db, errC = connectDB(d)

	if errC != nil {
		return
	}

	defer db.Close()

	if db == nil {
		Log.Error("What da heck bro, db is null ??")
		return
	}

	var query *sql.Rows
	var queryString string
	var err error

	// Créer la requête
	if condition != nil {
		query, err = db.Query("DELETE FROM " + table + " WHERE " + *condition)
		queryString = "DELETE FROM " + table + " WHERE " + *condition
		if err != nil {
			Log.Error("Error : ", err)
			Log.Debug(queryString)
			return
		}
	} else if *condition == "-1" {
		query, err = db.Query("DELETE FROM " + table)
		queryString = "DELETE FROM " + table
		if err != nil {
			Log.Error("Error : ", err)
			Log.Debug(queryString)
			return
		}
	}

	if err := query.Err(); err != nil {
		Log.Error("An error Occured : ", err)
		return
	}

	if len(debug) > 0 && debug[0] {
		Log.Debug(queryString)
	}
	Log.Infos("Deleting successful from " + table)
	return
}

//
// ------------------------------------------------------------------------------------------------ //
//

// transformQueryToMap Tranforme un tableau venant tout droit de la base de donnée avec les résultats
// Créer et renvoie un "objet json" construit dynamiquement avec les noms des colonnes et leur valeurs
func transformQueryToMap(query *sql.Rows) []map[string]interface{} {
	var result []map[string]interface{}

	// Pour tout les résultats
	for query.Next() {

		//Prend les colonnes
		columns, err := query.Columns()

		if err != nil {
			Log.Error("Impossible to retrieve column name")
			return nil
		}

		// Créer un slice pour stocker les vlaures
		values := make([]interface{}, len(columns))

		// Créer un pointeur de valeurs
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := query.Scan(pointers...); err != nil {
			Log.Error("Impossible to determine the pointer when creating the map")
			return nil
		}

		// Créer une sorte d'objet Json
		row := make(map[string]interface{})
		for i, name := range columns {
			switch v := values[i].(type) {
			case []byte:
				row[name] = string(v)
			default:
				row[name] = v
			}
		}

		result = append(result, row)
	}
	return result
}

//
// ------------------------------------------------------------------------------------------------ //
//

// checkData Regarde le bon format des données rentrées en paramètre et renvoie un booléen
// Utile pour les fonctions de base de donnée
func checkData(table string, column []string, values []string, condition *string) bool {

	if reflect.TypeOf(table) != reflect.TypeOf("") || table == "" {
		Log.Error("Need a tale name")
		return false
	}

	if reflect.TypeOf(column).Kind() != reflect.Slice || len(column) == 0 {
		Log.Error("Need a string array")
		return false
	}

	if reflect.TypeOf(values).Kind() != reflect.Slice || len(column) == 0 {
		Log.Error("Need a string array")
		return false
	}

	if condition != nil && reflect.TypeOf(*condition) != reflect.TypeOf("") {
		Log.Error("Need a string into the condition")
		return false
	}

	return true
}

//
// ------------------------------------------------------------------------------------------------ //
//

// ArrayToString Transforme un tableau en chaine de caractères
// Peut prendre un paramètre "noQuotes" qui évite de mettre des quotes lorsque que la fonction est utilisée pour transformer un tableau de nom de colonne en string
// Les quotes sont utiles pour les valeurs de ces colonnes pour un Insert pour les chaine de charactères
func ArrayToString(arr []string, noQuotes ...bool) string {
	if len(arr) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, s := range arr {

		// Vérifie si c'est un nombre avant pour mettre ou non des '
		_, err := strconv.Atoi(s)
		if err != nil && noQuotes == nil {
			sb.WriteString(`'` + s + `'`)
		} else {
			sb.WriteString(s)
		}

		if i < len(arr)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

//
// ------------------------------------------------------------------------------------------------ //
//

// ConcatColumnWithValues Utile pour les Update de BDD
// Plus particulièrement pour les "SET column=value, column=value, ..."
func ConcatColumnWithValues(columns []string, values []string) string {

	if len(columns) == 0 || len(values) == 0 {
		Log.Error("Plz columns and values string array must have at least one key each")
		return ""
	}

	if len(columns) != len(values) {
		Log.Error("Plz columns and values string array must have the same length")
		return ""
	}

	var sb strings.Builder
	for i, s := range values {

		// Vérifie si c'est un nombre avant pour mettre ou non des '
		_, err := strconv.Atoi(s)
		if err != nil {
			if s == "NULL" {
				sb.WriteString(columns[i] + `=` + s)
			} else {
				sb.WriteString(columns[i] + `='` + s + `'`)
			}
		} else {
			sb.WriteString(columns[i] + "=" + s)
		}

		if i < len(columns)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
