package main

import (
	"errors"
	"fmt"

	"log"
	"net/http"
	"net/smtp"

	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/HipolitoBautista/internal/models"
)

// written by: Hipolito, Michael, Jahmur, Dennis, Rene
// tested by: Hipolito, Michael, Jahmur, Dennis, Rene
// debugged by: Hipolito, Michael, Jahmur, Dennis, Rene

// Stores the form ID which needs to be shared between multiple handers
// shared data store
var dataStore = struct {
	sync.RWMutex
	data map[string]int64
}{data: make(map[string]int64)}

// Loads the UserUI (Dashboard)
func (app *application) UserUI(w http.ResponseWriter, r *http.Request) {
	//Checks for the user in the session,
	UserIDTemp := app.sessionsManager.GetInt(r.Context(), "authenticatedUserID")
	FD, err := app.form.IndividualLoadData(UserIDTemp)
	//Code to load the user's name on the top right of the Dashboard
	Username, _, _ := app.publicuser.UserData(UserIDTemp)
	accepted, unverified, denied, total, err := app.form.NormalStats(UserIDTemp)

	if err != nil {
		fmt.Println(err)
	}
	//Flash for notifications
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	fmt.Println(flash)

	//Passing our data into an instance of template data
	data := &templateData{
		Form:       FD,
		Flash:      flash,
		Accepted:   accepted,
		Unverified: unverified,
		Denied:     denied,
		Username:   Username,
		Total:      total,
	}

	//Parsing our tmpl file
	ts, err := template.ParseFiles("./ui/html/User.UI.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	//Passing  our instance of templatedata which now contains data into the tmpl file to display to the user
	err = ts.Execute(w, data)

	fmt.Println(err)
}

// Loads our admin TableView
func (app *application) TableView(w http.ResponseWriter, r *http.Request) {

	//Get all data to load
	FD, err := app.form.LoadData()

	if err != nil {
		fmt.Println(err)
	}

	//Getting the admins ID
	AdminIDTemp := app.sessionsManager.GetInt(r.Context(), "authenticatedAdminID")
	// getting the admin that's currently logged in
	Username, _, err := app.admin.AdminData(AdminIDTemp)
	//Getting the form stats
	accepted, unverified, denied, total, err := app.form.Stats()

	if err != nil {
		fmt.Println(err)
	}

	//Poping our flash
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	//Creating an instance of template data with data inside of it
	data := &templateData{
		Form:       FD,
		Flash:      flash,
		Accepted:   accepted,
		Unverified: unverified,
		Denied:     denied,
		Username:   Username,
		Total:      total,
	}

	//Paring out tmpl file
	ts, err := template.ParseFiles("./ui/html/Table.View.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	//passing our data into the tmpl file to show the user
	err = ts.Execute(w, data)
	fmt.Println(err)
}

func (app *application) ArchiveView(w http.ResponseWriter, r *http.Request) {
	//Loading data from archive table
	FD, err := app.form.LoadArchiveData()

	if err != nil {
		fmt.Println(err)
	}
	//Getting the total number of archived files
	total, err := app.archive.Stats()
	//Poping our flash
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	//Creating an instance of template data with data we need
	data := &templateData{
		Form:  FD,
		Flash: flash,
		Total: total,
	}
	//parsing our tmpl file
	ts, err := template.ParseFiles("./ui/html/Archive.View.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	//passing our data to the tmpl file to display to the user
	err = ts.Execute(w, data)
	fmt.Println(err)

}

// Loads a Form
func (app *application) LoadForm(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "./ui/html/Form.UI.tmpl", nil)

}

// Renders the Sign Up page template
func (app *application) SignUpPage(w http.ResponseWriter, r *http.Request) {
	//getting our flash to see if we have potential notifications to display to the user
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	//rendering our signup page
	RenderTemplate(w, "./ui/html/signup.UI.tmpl", data)

}

// Handles the data given to us in the signup page
func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	//Getting all data from the form
	r.ParseForm()
	Lname := r.PostForm.Get("Lname")
	Fname := r.PostForm.Get("Fname")
	name := strings.ReplaceAll(Fname, " ", "") + " " + strings.ReplaceAll(Lname, " ", "")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")

	//write the data to the table
	err := app.publicuser.Insert(name, email, password)
	//Checking for errors and redirecting after completion
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			app.sessionsManager.Put(r.Context(), "flash", "Email Already Registered")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
		} else {
			app.sessionsManager.Put(r.Context(), "flash", "Oops! Try again later")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
		}
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Signup was successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Renders the Login page with it's data
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	//flash for notifications to display to the user
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	//Rendering our Login page
	RenderTemplate(w, "./ui/html/Login.UI.tmpl", data)
}

// Handles the data given to us in the login page
func (app *application) loginSubmitted(w http.ResponseWriter, r *http.Request) {
	//Getting data submitted by the user
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")
	id, err := app.publicuser.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.sessionsManager.Put(r.Context(), "flash", "Invalid Credentials")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	// add user to the session cookie
	err = app.sessionsManager.RenewToken(r.Context())
	if err != nil {
		return
	}

	//add an authenticate user to the session
	app.sessionsManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/userui", http.StatusSeeOther)
}

// Loads admin login page
func (app *application) AdminLogin(w http.ResponseWriter, r *http.Request) {
	//flash for user notifications (alerts)
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	//Rendering our admin login page
	RenderTemplate(w, "./ui/html/AdminLogin.UI.tmpl", data)
}

// Handles the data submitted by adminloginpage when submit is clicked
func (app *application) AdminLoginSubmitted(w http.ResponseWriter, r *http.Request) {
	//Handing the submitted data
	r.ParseForm()
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")

	//Passing entered email and password for authentication of admin account
	id, err := app.admin.AuthenticateAdmin(email, password)
	//Error handing, setting a message for the user in flash and redirecting
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.sessionsManager.Put(r.Context(), "flash", "Invalid Credentials")
			http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
		}
		return
	}

	// add user to the session cookie
	err = app.sessionsManager.RenewToken(r.Context())
	if err != nil {
		fmt.Println(err)
		return
	}

	//add an authenticate admin with admin and user permissions in the session
	app.sessionsManager.Put(r.Context(), "authenticatedAdminID", id)
	app.sessionsManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/TableView", http.StatusSeeOther)
}

// Renders the Sign Up page template for admins
func (app *application) AdminSignUpPage(w http.ResponseWriter, r *http.Request) {
	//Checks for notifcation data stored inside flash to display to the user
	flash := app.sessionsManager.PopString(r.Context(), "flash")
	data := &templateData{
		Flash: flash,
	}
	//Rending our admin signup page
	RenderTemplate(w, "./ui/html/adminsignup.UI.tmpl", data)

}

// Handles the data given to us in the signup page for admins
func (app *application) Adminsignup(w http.ResponseWriter, r *http.Request) {
	//Gets admin data given for signup
	r.ParseForm()
	Lname := r.PostForm.Get("Lname")
	Fname := r.PostForm.Get("Fname")
	name := strings.ReplaceAll(Fname, " ", "") + " " + strings.ReplaceAll(Lname, " ", "")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("Password")
	auth := r.PostForm.Get("AuthenticationCode")

	//write the data to the table

	err := app.admin.InsertAdmin(name, email, password, auth)

	//Checking for errors in the process if the email is already in use
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			app.sessionsManager.Put(r.Context(), "flash", "Email Already Registered")
			http.Redirect(w, r, "/AdminSignup", http.StatusSeeOther)
		} else {
			app.sessionsManager.Put(r.Context(), "flash", "Oops! Try again later")
			http.Redirect(w, r, "/AdminSignup", http.StatusSeeOther)
		}
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Signup was successful")
		http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
	}

}

// Handles the when submit is pressed inside an affidavit
func (app *application) SubmittedAffidavit(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//get the form data

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return

	}
	//collecting values from form
	name := r.PostForm.Get("full-name")
	other_name := r.PostForm.Get("other-names")
	name_change, _ := strconv.ParseBool(r.PostForm.Get("name-changed"))
	previous_name := r.PostForm.Get("previous-names")
	reason_for_change := r.PostForm.Get("reason-for-change")
	social_num, _ := strconv.ParseInt(r.PostForm.Get("SS-number"), 10, 64)
	social_date, _ := time.Parse("2006-01-02", r.PostForm.Get("SS-issue-date"))
	social_country := r.PostForm.Get("SS-country-value")
	passport_num, _ := strconv.ParseInt(r.PostForm.Get("PP-number"), 10, 64)
	passport_date, _ := time.Parse("2006-01-02", r.PostForm.Get("PP-issue-date"))
	passport_country := r.PostForm.Get("PP-country-value")
	NHI_num, _ := strconv.ParseInt(r.PostForm.Get("NHI-number"), 10, 64)
	NHI_date, _ := time.Parse("2006-01-02", r.PostForm.Get("NHI-issue-date"))
	NHI_country := r.PostForm.Get("NHI-country-value")
	dob, _ := time.Parse("2006-01-02", r.PostForm.Get("dateOfBirth"))
	pob := r.PostForm.Get("placeOfBirth")
	nationality := r.PostForm.Get("nationality")
	acq_nationality := r.PostForm.Get("acq-nationality")
	spouse := r.PostForm.Get("spouse-name")
	address := r.PostForm.Get("AF-address")
	phone, _ := strconv.ParseInt(r.PostForm.Get("AF-number"), 10, 64)
	fax, _ := strconv.ParseInt(r.PostForm.Get("Fax-Number"), 10, 64)
	email := r.PostForm.Get("email-address")

	//Gets the user ID of the user submitting the form
	UserID := app.sessionsManager.GetInt(r.Context(), "authenticatedUserID")
	//create instance of form model and passes the form data to it
	data := &models.Form{
		User_id:                  int32(UserID),
		Form_status:              "unverified",
		Archive_status:           false,
		Affiant_full_name:        name,
		Other_names:              other_name,
		Name_change_status:       name_change,
		Previous_name:            previous_name,
		Reason_for_Change:        reason_for_change,
		Social_security_num:      social_num,
		Social_security_date:     social_date,
		Social_security_country:  social_country,
		Passport_number:          passport_num,
		Passport_date:            passport_date,
		Passport_country:         passport_country,
		NHI_number:               NHI_num,
		NHI_date:                 NHI_date,
		NHI_country:              NHI_country,
		Dob:                      dob,
		Place_of_birth:           pob,
		Nationality:              nationality,
		Acquired_nationality:     acq_nationality,
		Spouse_name:              spouse,
		Affiants_address:         address,
		Residencial_phone_number: phone,
		Residenceial_fax_num:     fax,
		Residencial_email:        email,
	}
	//inserting data into DB using the insert method
	form_id, err := app.form.Insert(data)
	fmt.Println(form_id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//redirects to the dashboard upon completion
	http.Redirect(w, r, "/userui", http.StatusSeeOther)

}

// Loads the form to be updated by the user
func (app *application) updateForm(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	ID, _ := strconv.ParseInt(r.PostForm.Get("edit"), 10, 64)
	//form ID to load when wanting to edit the form (READING)
	form_id := ID
	FD, err := app.form.Read(form_id)
	CD, err := app.comments.GetComments(form_id)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	//Checking to see if the user is an admin, giving them access to admin buttons etc
	var AdminPerms string
	if app.sessionsManager.Exists(r.Context(), "authenticatedAdminID") {
		AdminPerms = "true"

	} else {
		AdminPerms = "false"
	}

	//an instance of templateData
	data := &templateData{
		Form:           FD,
		UserPermStatus: AdminPerms,
		Comments:       CD,
	} //this allows us to pass in mutliple data into the template

	//parsing out files
	ts, err := template.ParseFiles("./ui/html/form.show.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	//getting the form ID
	dataStore.Lock()
	dataStore.data["key"] = int64(form_id)
	dataStore.Unlock()

	//handing our data to the tmpl file
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

}

// Handles the POST from the affidavit update form
func (app *application) updateFormQuery(w http.ResponseWriter, r *http.Request) {
	//Code to UPDATE a form

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//get the form data

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return

	}
	//collecting values from form
	//form_id, _ := strconv.ParseInt(r.PostForm.Get("form-id"), 10, 64)
	name := r.PostForm.Get("full-name")
	other_name := r.PostForm.Get("other-names")
	name_change, _ := strconv.ParseBool(r.PostForm.Get("name-changed"))
	previous_name := r.PostForm.Get("previous-names")
	reason_for_change := r.PostForm.Get("reason-for-change")
	social_num, _ := strconv.ParseInt(r.PostForm.Get("SS-number"), 10, 64)
	social_date, _ := time.Parse("2006-01-02", r.PostForm.Get("SS-issue-date"))
	social_country := r.PostForm.Get("SS-country-value")
	passport_num, _ := strconv.ParseInt(r.PostForm.Get("PP-number"), 10, 64)
	passport_date, _ := time.Parse("2006-01-02", r.PostForm.Get("PP-issue-date"))
	passport_country := r.PostForm.Get("PP-country-value")
	NHI_num, _ := strconv.ParseInt(r.PostForm.Get("NHI-number"), 10, 64)
	NHI_date, _ := time.Parse("2006-01-02", r.PostForm.Get("NHI-issue-date"))
	NHI_country := r.PostForm.Get("NHI-country-value")
	dob, _ := time.Parse("2006-01-02", r.PostForm.Get("dateOfBirth"))
	pob := r.PostForm.Get("placeOfBirth")
	nationality := r.PostForm.Get("nationality")
	acq_nationality := r.PostForm.Get("acq-nationality")
	spouse := r.PostForm.Get("spouse-name")
	address := r.PostForm.Get("AF-address")
	phone, _ := strconv.ParseInt(r.PostForm.Get("AF-number"), 10, 64)
	fax, _ := strconv.ParseInt(r.PostForm.Get("Fax-Number"), 10, 64)
	email := r.PostForm.Get("email-address")
	dataStore.RLock()
	form_id := dataStore.data["key"]
	dataStore.RUnlock()

	//Creating an instance of template data to pass data to the tmpl file

	data := &models.Form{
		Form_id:                  form_id,
		Form_status:              "unverified",
		Archive_status:           false,
		Affiant_full_name:        name,
		Other_names:              other_name,
		Name_change_status:       name_change,
		Previous_name:            previous_name,
		Reason_for_Change:        reason_for_change,
		Social_security_num:      social_num,
		Social_security_date:     social_date,
		Social_security_country:  social_country,
		Passport_number:          passport_num,
		Passport_date:            passport_date,
		Passport_country:         passport_country,
		NHI_number:               NHI_num,
		NHI_date:                 NHI_date,
		NHI_country:              NHI_country,
		Dob:                      dob,
		Place_of_birth:           pob,
		Nationality:              nationality,
		Acquired_nationality:     acq_nationality,
		Spouse_name:              spouse,
		Affiants_address:         address,
		Residencial_phone_number: phone,
		Residenceial_fax_num:     fax,
		Residencial_email:        email,
	}

	//Calling our Update function which to update the form in the DB
	_, err = app.form.Update(data)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//Setting a flash (notification) for the user
	app.sessionsManager.Put(r.Context(), "flash", "Form edited Successfully")

	//Getting comments to store in the DB from form if the user is an admin
	if app.isAuthenticatedAdmin(r) == true {

		removeThis := "\\r+"
		//Trimming all comment box spaces
		comment1 := strings.TrimSpace(strings.Replace(r.PostForm.Get("1"), removeThis, "", -1))
		comment2 := strings.TrimSpace(strings.Replace(r.PostForm.Get("2"), removeThis, "", -1))
		comment3 := strings.TrimSpace(strings.Replace(r.PostForm.Get("3"), removeThis, "", -1))
		comment35 := strings.TrimSpace(strings.Replace(r.PostForm.Get("3.5"), removeThis, "", -1))
		comment4 := strings.TrimSpace(strings.Replace(r.PostForm.Get("4"), removeThis, "", -1))
		comment5 := strings.TrimSpace(strings.Replace(r.PostForm.Get("5"), removeThis, "", -1))
		comment6 := strings.TrimSpace(strings.Replace(r.PostForm.Get("6"), removeThis, "", -1))
		comment7 := strings.TrimSpace(strings.Replace(r.PostForm.Get("7"), removeThis, "", -1))
		comment8 := strings.TrimSpace(strings.Replace(r.PostForm.Get("8"), removeThis, "", -1))
		comment9 := strings.TrimSpace(strings.Replace(r.PostForm.Get("9"), removeThis, "", -1))
		comment10 := strings.TrimSpace(strings.Replace(r.PostForm.Get("10"), removeThis, "", -1))

		//If  all comments are empty do not insert anything, else insert the comments
		var ShouldInsert bool
		if comment1 != "" || comment2 != "" || comment3 != "" || comment35 != "" || comment4 != "" || comment5 != "" || comment6 != "" || comment7 != "" || comment8 != "" || comment9 != "" || comment10 != "" {
			ShouldInsert = true
		}

		//Getting ID of admin that is logged in currently
		AdminIDTemp := app.sessionsManager.GetInt(r.Context(), "authenticatedAdminID")
		AdminID := int64(AdminIDTemp)

		//If comments are not empty we should insert the comments by creating an instance of our models.comments. This instance stores our data to be passed.
		if ShouldInsert == true {
			commentData := &models.Comments{
				Admin_id:  AdminID,
				Form_id:   form_id,
				Comment1:  comment1,
				Comment2:  comment2,
				Comment3:  comment3,
				Comment35: comment35,
				Comment4:  comment4,
				Comment5:  comment5,
				Comment6:  comment6,
				Comment7:  comment7,
				Comment8:  comment8,
				Comment9:  comment9,
				Comment10: comment10,
			}

			fmt.Println(comment1)
			//Checking to see if the form exists to add comments to it
			if app.comments.Exists(form_id) == true {
				//If the form already exists we UPDATE the comments
				app.comments.UpdateComments(commentData)
			} else {
				//If it doesnt exist yet we INSERT the comment.
				app.comments.InsertComment(commentData)
			}
		}
		//Redirecting upon completion
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/userui", http.StatusSeeOther)

	}

}

type json_data struct {
	formid string
}

// Handles form Deletion
func (app *application) deleteForm(w http.ResponseWriter, r *http.Request) {
	//Code to delete a form
	//Gets formID of form to be deleted
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("del"), 10, 64)

	if err != nil {
		panic(err)
	}

	//Passes it to the Delete function which runs the SQL to delete from DB
	err = app.form.Delete(ID)

	//Error handing, setting notifications and redirecting upon completion
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Deletion Failed, Try Again Later")
		http.Redirect(w, r, "/userui", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Successful Deletion")
		http.Redirect(w, r, "/userui", http.StatusSeeOther)
	}

}

// changes affidavit status to accepted
func (app *application) affidavitAccept(w http.ResponseWriter, r *http.Request) {
	//gets the form ID of the form to be changed to accepted
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("accept"), 10, 64)

	if err != nil {
		panic(err)
	}
	//Passes ID to the AcceptForm function which runs the SQL to change that field within the table to status for this form.
	err = app.form.AcceptForm(ID)

	//Error handling, setting notifications and redirecting once task is completed
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Action Failed, Try Again Later")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "form has been accepted")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	}

}

// changes form status to denied
func (app *application) affidavitDeny(w http.ResponseWriter, r *http.Request) {
	//Getting the form ID of the form to be denied
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("deny"), 10, 64)
	IDString := strconv.FormatInt(ID, 10)

	if err != nil {
		panic(err)
	}
	//Gets data from the form owner
	Owner, err := app.form.FormOwner(ID)
	if err != nil {
		fmt.Println(err)
	}
	//Gets the email assocaited with the owner of this form
	_, email, err := app.publicuser.UserData(Owner)
	if err != nil {
		fmt.Println(err)
	}

	//passing form Id into the DenyForm function which changes the forms status to denied in the DB (runs SQL to do this)
	err = app.form.DenyForm(ID)
	//Error handling , setting notifications and Redirecting
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Action Failed, Try Again Later")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Form has been denied ")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)

	}

	//Automatically sends email to inform the user the form has been denied

	from := "osippproject@gmail.com"
	password := "pzkrxnbaznlannih"
	to := email

	//giving it the server to authenticate details.
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	smtpServer := "smtp.gmail.com:587"
	conn, err := smtp.Dial(smtpServer)
	if err != nil {
		panic(err)
	}
	//closes connection once email is sent
	defer conn.Close()

	// message to be emailed
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Form Has been denied\r\n" +
		"\r\n" +
		"Hello,\r\n" +
		"We regret to inform you that your form (" + IDString + ") has been denied. \r\n" +
		"Your form may have been denied for several reasons.\r\n" +
		"In order to fix this issue log into your OSIPP account and open the denied form.\r\n" +
		"Within the form change the fields according to the comments our OSIPP Officers have left. Once the changes have been made resubmit the form\r\n")

	// Send the email.
	if err := smtp.SendMail(smtpServer, auth, from, []string{to}, msg); err != nil {
		panic(err)
	}

}

// moves form to the archived table
func (app *application) ArchiveForm(w http.ResponseWriter, r *http.Request) {
	//Get ID of the form to be archived
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("archive"), 10, 64)

	if err != nil {
		panic(err)
	}

	//Passes the form ID to ArchiveForm function which changes the form to the Archive table
	err = app.archive.ArchiveForm(ID)

	//Notification setting, Error handling, and Redirecting
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Archive Failed, Try Again Later")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Form has been Archived ")
		http.Redirect(w, r, "/TableView", http.StatusSeeOther)
	}

}

// Moves form from archive to formarchive
func (app *application) UnarchiveForm(w http.ResponseWriter, r *http.Request) {
	//Getting ID of the form to be unarchived
	err := r.ParseForm()
	ID, _ := strconv.ParseInt(r.PostForm.Get("unarchive"), 10, 64)

	if err != nil {
		panic(err)
	}

	//Passes the ID to the unarchive function which moves the form from the archive table to the form table
	err = app.archive.UnArchiveForm(ID)
	fmt.Println(err)
	//setting notifications, redirecting and error handling
	if err != nil {
		app.sessionsManager.Put(r.Context(), "flash", "Archive Failed, Try Again Later")
		http.Redirect(w, r, "/ArchiveView", http.StatusSeeOther)
	} else {
		app.sessionsManager.Put(r.Context(), "flash", "Form has been Archived ")
		http.Redirect(w, r, "/ArchiveView", http.StatusSeeOther)
	}
}

// logs out the user
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {

	//If the user is an admin
	if app.isAuthenticatedAdmin(r) == true {
		//Remove admin from the session
		app.sessionsManager.PopString(r.Context(), "authenticatedUserID")
		app.sessionsManager.PopString(r.Context(), "authenticatedAdminID")
		app.sessionsManager.Put(r.Context(), "flash", "Successfully Logged out ")
		http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
	} else {
		//else the user is not an admin, remove user from session
		app.sessionsManager.PopString(r.Context(), "authenticatedUserID")
		app.sessionsManager.Put(r.Context(), "flash", "Successfully Logged out ")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}
