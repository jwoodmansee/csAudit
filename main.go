package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Audit : Model audit types
// type Audit struct {
// 	AuditID   string
// 	Name      string
// 	AccountID string
// }

// AData : Model for future database call
type AData struct {
	// Auditor        string
	// AuditorID      string
	// ReciepentID    string
	AuditID        string
	Name           string
	CatagoryOrder  string
	CatagoryID     string
	CatagoryName   string
	CatagoryLabel  string
	QuestionOrder  string
	QuestionID     string
	QuestionText   string
	PointsEarned   string
	PointsPossible string
	IsOptional     string
}

// Audit . . .
type Audit struct {
	AuditID    string
	Name       string
	AccountID  string
	Catagories []Catagory
	Questions  []Question
}

//Catagory . . .
type Catagory struct {
	Order string
	ID    string
	Label string
}

// Question . . .
type Question struct {
	Order string
	ID    string
	Text  string
}

func getIndex(arr []Audit, ID string) int {
	for k, v := range arr {
		if v.AuditID == ID {
			return k
		}
	}
	return -1
}

// GetAudit handle request for what audit a person can do
var GetAudit = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	audits, err := getAudit(id)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	a, err := json.Marshal(audits)
	if err != nil {
		fmt.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(a)

})

func getAudit(id string) ([]*Audit, error) {

	



	audits := []*Audit{}
	for rows.Next() {
		a := &Audit{}

		if err = rows.Scan(
			&a.AuditID,
			&a.Name,
			&a.AccountID); err != nil {
			return nil, err
		}
		audits = append(audits, a)
	}
	return audits, nil
}

// GetAuditForm handle request for survey data
var GetAuditForm = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")

	decoder := json.NewDecoder(r.Body)
	var a Audit
	err := decoder.Decode(&a)
	if err != nil {
		panic(err)
	}

	var body string = a.AuditID

	fData, err := getAuditInfo(id, body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	form, err := json.Marshal(fData)
	if err != nil {
		fmt.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(form)
})

func getAuditInfo(id string, body string) ([]*AData, error) {
	

	form := []*AData{}
	for rows.Next() {
		f := &AData{}

		if err = rows.Scan(
			&f.AuditID,
			&f.Name,
			&f.CatagoryOrder,
			&f.CatagoryID,
			&f.CatagoryName,
			&f.CatagoryLabel,
			&f.QuestionOrder,
			&f.QuestionID,
			&f.QuestionText,
			&f.PointsPossible,
			&f.IsOptional); err != nil {
			return nil, err
		}
		form = append(form, f)
	}
	return form, nil
}

// UpdateAudit handles request for update audits
var UpdateAudit = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")

	success := putAudit(id, r)

	if success {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotModified)
	}
})

func putAudit(id string, r *http.Request) bool {
	success := false

	db := data.GetDB()
	if db == nil {
		return success // false
	}
	defer db.Close()

	bytes, _ := ioutil.ReadAll(r.Body)
	body := string(bytes)
	fmt.Println(body)

	row, err := db.Query("set nocount on; exec [spcAuditCreate] ?, ?", id, body)
	if err != nil {
		log.Println("Query failed: ", err.Error())
		return success // false
	}
	defer row.Close()

	cols, err := row.Columns()
	if err != nil {
		return success // false
	}
	if cols == nil {
		return false // false
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	for row.Next() {
		err = row.Scan(vals...)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		replyFromDB := data.ParseValue(vals[0].(*interface{}))
		success, err = strconv.ParseBool(replyFromDB)
		if err != nil {
			log.Println(err.Error())
			return success
		}
	}

	if row.Err() != nil {
		return false
	}

	return success
}

func getAuditInfo2(id string, body string) []byte {
	db := data.GetDB()
	if db == nil {
		return nil
	}
	defer db.Close()

	rows, err := db.Query("set nocount on; exec [spcAuditGet] ?, ?", id, body)
	for err != nil {
		log.Println(err.Error())
		return nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if cols == nil {
		log.Println(cols)
		return nil
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	var q []AData

	for rows.Next() {
		row := AData{
			AuditID:        "",
			Name:           "",
			CatagoryOrder:  "",
			CatagoryID:     "",
			CatagoryName:   "",
			CatagoryLabel:  "",
			QuestionID:     "",
			QuestionOrder:  "",
			QuestionText:   "",
			PointsEarned:   "",
			PointsPossible: "",
			IsOptional:     "",
		}
	}
}

func main() {

}
