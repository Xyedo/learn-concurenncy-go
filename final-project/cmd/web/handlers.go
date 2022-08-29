package main

import (
	"errors"
	"final-project/data"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}
func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	err := app.Session.RenewToken(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	isValidPass, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !isValidPass {
		msg := Message{
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Invalid login attempt!",
		}

		app.sendEmail(msg)
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "Successful login!")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	err := app.Session.Destroy(r.Context())
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = app.Session.RenewToken(r.Context())
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}
func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	email := r.Form.Get("email")
	user, _ := app.Models.User.GetByEmail(email)
	if user != nil {
		app.Session.Put(r.Context(), "error", "Email already in use")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	newUser := data.User{
		Email:     email,
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}
	_, err = newUser.Insert(newUser)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	url := fmt.Sprintf("http://localhost/activate?email=%s", newUser.Email)
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       newUser.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.sendEmail(msg)

	app.Session.Put(r.Context(), "flash", "Confirmation email sent, Check your email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	testURL := fmt.Sprintf("http://localhost%s", url)
	if okay := VerifyToken(testURL); !okay {
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u.Active = 1
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update User.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "flash", "Account activated, You can log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
	dataMap := make(map[string]any)

	dataMap["plans"] = plans
	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
	plan, err := app.Models.Plan.GetOne(id)
	if err != nil {
		app.Session.Put(r.Context(), "error", "cannot get the plans")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}
	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Log in first")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		invoice, err := app.getInvoice(&user, plan)
		if err != nil {
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your invoice",
			Data:     invoice,
			Template: "invoice",
		}

		app.sendEmail(msg)

	}()
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		pdf := app.gerenateManual(&user, plan)
		err := pdf.OutputFileAndClose(fmt.Sprintf("./tmp/%d_manual.pdf", user.ID))
		if err != nil {
			app.ErrorChan <- err
			return
		}
		msg := Message{
			To:      user.Email,
			Subject: "Your manual",
			Data:    "Your user manual is attached",
			AttachmentMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("./tmp/%d_manual.pdf", user.ID),
			},
		}
		app.sendEmail(msg)
		app.ErrorChan <- errors.New("some custom errors")
	}()
	if err = app.Models.Plan.SubscribeUserToPlan(user, *plan); err != nil {
		app.Session.Put(r.Context(), "error", "Error subscribing to plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}
	u, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error getting user from database ")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "user", u)

	app.Session.Put(r.Context(), "flash", "Subscribed!")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}
func (app *Config) getInvoice(u *data.User, plan *data.Plan) (string, error) {
	return plan.PlanAmountFormatted, nil
}

func (app *Config) gerenateManual(u *data.User, plan *data.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")

	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()
	time.Sleep(5 * time.Second)
	t := importer.ImportPage(pdf, "./pdf/manual.pdf", 1, "/MediaBox")
	pdf.AddPage()
	importer.UseImportedTemplate(pdf, t, 0, 0, 251.9, 0)
	pdf.SetX(75)
	pdf.SetY(150)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName, u.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)
	return pdf
}
