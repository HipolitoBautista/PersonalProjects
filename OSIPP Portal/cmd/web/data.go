package main

import "github.com/HipolitoBautista/internal/models"

  // written by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // tested by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // debugged by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // etc. 
  

  //This is our templdateData structure which is used to hand in multiple pieces of data to the tmpl files 
type templateData struct {
	Form           []*models.Form
	Flash          string
	UserPermStatus string
	Comments       []*models.Comments
	Accepted       string
	Denied         string
	Unverified     string
	Username       string
	Total          string
}
