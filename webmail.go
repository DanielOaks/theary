package main

/*
 * This file is part of theary.
 *
 * theary is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * theary is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/gorilla/mux"
)

type dataTable struct {
	//Echo   				int				`json:"sEcho"`
	TotalRecords        int        `json:"iTotalRecords"`
	TotalDisplayRecords int        `json:"iTotalDisplayRecords"`
	Rows                [][]string `json:"aaData"`
}

type EmailTable struct {
	Cells []string `json:",string"`
}

// homeView displays a minimalist webmail client  (renders template)
func homeView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p := Page{Title: "Test"}
	err := tmpl.Execute(w, p)
	checkHttpError(err, w)
}

// listMailsWS list the received e-mails from db (returns JSON)
func listMailsWS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	recipient := vars["recipient"]
	emails := dbEmails.Use(recipient)

	var mailRecords dataTable
	//mailRecords.Echo = 3

	queryStr := `"all"`
	var query interface{}
	var record map[string]interface{}
	json.Unmarshal([]byte(queryStr), &query)
	queryResult := make(map[int]struct{})
	err := db.EvalQuery(query, emails, &queryResult)
	checkHttpError(err, w)

	for id := range queryResult {
		record, err = emails.Read(id)
		mailRecords.TotalRecords++
		mailRecords.TotalDisplayRecords++
		//sow := []string{strconv.FormatInt(int64(id), 10),
		row := []string{strconv.Itoa(id),
			record["timestamp"].(string),
			record["from"].(string),
			record["subject"].(string),
			record["address"].(string)}
		//var rowEntry EmailTable
		//rowEntry.Cells = row
		mailRecords.Rows = append(mailRecords.Rows, row)
	}
	jsonString, err := json.Marshal(mailRecords)
	checkHttpError(err, w)
	fmt.Fprintf(w, "%s", jsonString)
}

// getMailWS returns the e-mails details from DB (returns JSON)
func getMailWS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	recipient := vars["recipient"]
	var id int
	//id, err := strconv.ParseInt(vars["id"], 10, 64)
	id, err := strconv.Atoi(vars["id"])
	checkHttpError(err, w)
	emails := dbEmails.Use(recipient)
	var record map[string]interface{}
	record, err = emails.Read(id)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(record["data"].(string))
	checkHttpError(err, w)
}

// checkRecipientWS checks if a recipient exists or not (returns JSON)
func checkRecipientWS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(r)
	recipient := vars["recipient"]
	encoder := json.NewEncoder(w)
	err := encoder.Encode(existsIndB(recipient))
	checkHttpError(err, w)
}
