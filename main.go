package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"vdicalc/auth"
	"vdicalc/calculations"
	"vdicalc/citrixcloud/trust"
	"vdicalc/citrixcloud/was"
	c "vdicalc/config"
	"vdicalc/functions"
	"vdicalc/mysql"
	"vdicalc/secretmanager"

	"github.com/spf13/viper"
	"google.golang.org/api/oauth2/v2"
)

var configuration c.Configurations
var db *sql.DB
var tokeninfo *oauth2.Tokeninfo // Maintain user's token authentication during sessions

func init() {

	viper.SetConfigName("config")   // Set the file name of the configurations file
	viper.AddConfigPath("./config") // Set the path to look for the configurations file
	viper.AutomaticEnv()            // Enable VIPER to read Environment Variables
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

}

func main() {

	/* Determine port for HTTP service. */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	/* Inititalize DB connection */
	db = mysql.DBInit()

	/* This function call up a function (functions/Dataload) to dynamicaly populate the dataset for the HTML files */
	FullData := functions.DataLoad(configuration)

	switch r.Method {
	case "GET":

		/* Determine the URI path to de taken */
		switch r.URL.Path {

		case "/":

			/* Determine the backend service address to be passsed to index.html js for authentication.
			AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
			FullData["authaddress"] = os.Getenv("AUTH_ADDRESS")

			/* This is the template execution for 'index' */
			functions.ExecuteTemplate(w, "index.html", FullData)

		case "/ccmetrix":

			/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
			if tokeninfo == nil {

				/* This  function return error codes to html */
				FullData["errorresults"] = functions.ReturnError("Warning", "You must be Signed In")

				/* Determine the backend service address to be passsed to index.html js for authentication.
				AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
				FullData["authaddress"] = os.Getenv("AUTH_ADDRESS")

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", FullData)

				break

			} else {

				/* Retrieve Citrix clientSecret from Google Secret Manager */
				secret := secretmanager.GetSecret("893974452758", tokeninfo.UserId)

				switch secret {
				case nil:

					/* Determine the backend service address to be passsed to index.html js for authentication.
					AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
					FullData["authaddress"] = os.Getenv("AUTH_ADDRESS")

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "ccmetrix_signup.html", FullData)

				default:

					var score int                      /* variable for score */
					var ranking map[string]interface{} /* variable for ranking */
					var startTime int64                /* variable for RequestUserExperienceTrend startTime */
					var endTime int64                  /* variable RequestUserExperienceTrend endTime */
					var clients trust.Clients          /* variable for Citrix Cloud bearer token */

					/* Request token from Citrix Cloud with credentials obtained from Google Secrets Manager */
					if clients.Token == "" {
						var err error
						clients, err = trust.RequestToken(secret["customerID"], secret["clientID"], secret["clientSecret"])

						/* Delete clientSecret and redirect if token is not valid */
						if err != nil {
							/* Delete Citrix clientSecret from Google Secret Manager */
							secretmanager.DeleteSecret("893974452758", tokeninfo.UserId)

							/* Redirect to / */
							http.Redirect(w, r, "/ccmetrix", http.StatusSeeOther)
						}
					}

					/* If IS_PROD environment variable is set, it contains
					True or False. If IS_PROD is not set, DEV URI is used. */
					var isDev bool
					if os.Getenv("IS_DEV") == "true" {
						isDev = true
						/* If "/ccmetrix?dev" used enters dev mode */
					} else if r.URL.RawQuery == "dev" {
						isDev = true
					}

					/* Retrieve agregate data for last 12 hours. */
					startTime = functions.TimetoEpoch(time.Now(), -720)
					endTime = functions.TimetoEpoch(time.Now(), 0)
					data := was.RequestUserExperienceTrend(clients, secret["customerID"], 30, startTime, endTime, "12h", "ALL", false, isDev)

					/* If data.TotalUsers == 0 results, try agregate data for last 1 week */
					if data.TotalUsers == 0 {
						startTime = functions.TimetoEpoch(time.Now(), -10080)
						endTime = functions.TimetoEpoch(time.Now(), 0)
						data = was.RequestUserExperienceTrend(clients, secret["customerID"], 360, startTime, endTime, "1w", "ALL", false, isDev)
					}

					/* If data.TotalUsers == 0 results, try agregate data for last 1 month */
					if data.TotalUsers == 0 {
						startTime = functions.TimetoEpoch(time.Now(), -43800)
						endTime = functions.TimetoEpoch(time.Now(), 0)
						data = was.RequestUserExperienceTrend(clients, secret["customerID"], 1440, startTime, endTime, "1m", "ALL", false, isDev)
					}

					/* Cleanup Citrix clientSecret */
					secret = nil    /* Cleanup Goole Secret */
					clients.Reset() /* Cleanup Citrix Cloud tokens */

					/* If data.TotalUsers == 0 results, assign score 0 or calculate score and ranking */
					if data.TotalUsers == 0 {
						/* Assign Ranking and Score 0 */
						score = 0
					} else {
						/* Calculate score and ranking */
						score = was.CalculateScore(data.TotalUsers, data.Items[0].Value, data.Items[1].Value, data.Items[2].Value)
						ranking = mysql.LoadScoreRanking(db, tokeninfo.UserId)
					}

					/* Test if user already exist and if not add to the table or update*/
					if mysql.QueryCCMetrixUser(db, tokeninfo.UserId) == false {
						mysql.SaveCCMetrixTransaction(db, tokeninfo.UserId, functions.GetIP(r), score)
					} else {
						mysql.UpdateCCMetrixTransaction(db, tokeninfo.UserId, functions.GetIP(r), score)
					}

					/* Determine the backend service address to be passsed to index.html js for authentication.
					AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
					FullData["authaddress"] = os.Getenv("AUTH_ADDRESS")
					FullData["score"] = score
					FullData["ranking"] = ranking["pos"]

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "ccmetrix.html", FullData)

				}

			}

		}

	case "POST":

		/* Determine the URI path to de taken */
		switch r.URL.Path {

		case "/tokensignoff":

			/* The index.html signoff JS triggers a http POST on /tokensignoff */
			tokeninfo = nil /* Clean up Google tokens */

		case "/tokensignin":

			/* The index.html signin JS triggers a http POST on /tokensignin and the user's token is associated to tokeninfo for verification*/
			if tokeninfo == nil {

				/* This function will verify user's token validity with Google GCP service */
				tokeninfo, _ = auth.VerifyIDToken(strings.TrimPrefix(r.URL.RawQuery, "id_token="))

				/* Test if user exist in mySQL vdicalc.users table, and if not add user to the table*/
				if mysql.QueryUser(db, tokeninfo.UserId) == false {

					/* This function executes the SQL estatement on Google SQL Run database */
					mysql.CreateUser(db, tokeninfo.UserId, tokeninfo.Email)

				}
			}

		case "/":

			/* This function reads and parse the html form */
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			r.ParseForm()

			/* This function loops throught the HTML form input fields to collect and store values during profile changes.
			Since values are re-loaded from config.yml existing form values are stored on *selected keys */
			for key, values := range r.Form {

				newKey := key + "selected"
				FullData[newKey] = values[0]

			}

			/* This function uses a hidden field 'submitselect' in each HTML template to detect the actions triggered by users.
			HTML action must include 'document.getElementById('submitselect').value='about';this.form.submit()' */
			switch r.PostFormValue("submitselect") {

			case "about":

				/* This is the template execution for 'about' */
				functions.ExecuteTemplate(w, "about.html", FullData)

			case "statistics":

				/* This function execute statistcs stored procedure */
				var data map[string]interface{} = mysql.LoadStatistics(db)

				/* This map variable stored results from LoadStatistics.*/
				for key := range data {

					FullData[key] = data[key]

				}

				/* This is the template execution for 'statistics' */
				functions.ExecuteTemplate(w, "statistics.html", FullData)

			case "ccmetrix":

				/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
				if tokeninfo == nil {

					/* This  function return error codes to html */
					FullData["errorresults"] = functions.ReturnError("Warning", "You must be Signed In")

					/* Determine the backend service address to be passsed to index.html js for authentication.
					AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
					FullData["authaddress"] = os.Getenv("AUTH_ADDRESS")

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "index.html", FullData)

					break

				} else {

					data := r.Header.Get("Origin") + "/ccmetrix"
					FullData["url"] = data
					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "ccmetrix_load.html", FullData)
				}

			case "ccmetrixsignupsubmit":

				/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
				if tokeninfo == nil {

					/* This  function return error codes to html */
					FullData["errorresults"] = functions.ReturnError("Warning", "You must be Signed In")

					/* Determine the backend service address to be passsed to index.html js for authentication.
					AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
					FullData["authaddress"] = os.Getenv("AUTH_ADDRESS")

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "ccmetrix_signup.html", FullData)

					break

				} else {

					/* This function reads and parse the html form */
					if err := r.ParseForm(); err != nil {
						fmt.Fprintf(w, "ParseForm() err: %v", err)
						return
					}
					r.ParseForm()

					/* Create Secret with Citrix Cloud API authorization */
					secretmanager.CreateSecret("893974452758", tokeninfo.UserId, r.PostFormValue("customerID"), r.PostFormValue("clientID"), r.PostFormValue("clientSecret"))

					/* Redirect to / */
					http.Redirect(w, r, "/ccmetrix", http.StatusSeeOther)

				}

			case "ccmetrixreset":

				/* Delete Citrix clientSecret from Google Secret Manager */
				secretmanager.DeleteSecret("893974452758", tokeninfo.UserId)

				/* Redirect to / */
				http.Redirect(w, r, "/ccmetrix", http.StatusSeeOther)

			case "back", "guide":

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", FullData)

			case "load":

				/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
				if tokeninfo == nil {

					/* This  function return error codes to html */
					FullData["errorresults"] = functions.ReturnError("Warning", "You must be Signed In")

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "index.html", FullData)

					break

				} else {

					/* This function load saved configurations */
					var data interface{} = mysql.LoadUserSaves(db, tokeninfo.UserId)
					FullData["usersaves"] = data

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "index.html", FullData)

				}

			case "usersaves":

				/* This map variable store results from LoadSaveByID. This function retrieve data from an existing configuration. */
				var data map[string]interface{} = mysql.LoadSaveByID(db, r.PostFormValue("usersaves"))

				for key := range data {

					FullData[key] = data[key]

				}

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", FullData)

			case "save":

				/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
				if tokeninfo == nil {

					/* This  function return error codes to html */
					FullData["errorresults"] = functions.ReturnError("Warning", "You must be Signed In")

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "index.html", FullData)

					break

				} else {

					/* This if function test if the savename provided is not empty and return a warning is it is */
					if len(r.PostFormValue("savename")) != 0 {

						/* This function save a configurationto into the database*/
						mysql.SaveConfiguration(db, tokeninfo.UserId, r.PostFormValue("savename"), FullData)

					} else {

						/* This  function return error codes to html */
						FullData["errorresults"] = functions.ReturnError("Warning", "Your save must have a name")

					}

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "index.html", FullData)
				}

			case "vmprofile":
				/* This is the execution case for profile change */

				var key c.VMConfigurations

				switch r.PostFormValue("vmprofile") {
				case "1":

					/* This is the profile execution for Task user */
					key = configuration.VMProfile01

				case "2":

					/* This is the profile execution for Office user */
					key = configuration.VMProfile02

				case "3":

					/* This is the profile execution for Knowledge user */
					key = configuration.VMProfile03

				case "4":

					/* This is the profile execution for Power user */
					key = configuration.VMProfile04
				}

				FullData["vmvcpucountselected"] = key.Vcpucountselected
				FullData["vmvcpumhzselected"] = key.Vcpumhz
				FullData["vmpercorecountselected"] = key.Vmpercorecountselected
				FullData["vmmemorysizeselected"] = key.Memorysize
				FullData["vmdisplaycountselected"] = key.Displaycountselected
				FullData["vmdisplayresolutionselected"] = key.Displayresolutionselected
				FullData["vmvideoramselected"] = key.Videoramselected
				FullData["vmdisksizeselected"] = key.Disksizeselected
				FullData["vmiopscountselected"] = key.Iopscountselected
				FullData["vmiopsreadratioselected"] = key.Iopsreadratioselected
				FullData["vmiopsbootcountselected"] = key.Iopsbootcountselected
				FullData["vmiopsbootreadratioselected"] = key.Iopsbootreadratioselected
				FullData["vmclonesizerefreshrateselected"] = key.Clonesizerefreshrateselected

				/* Fallthrough ensures 'update' is always executed after 'vmprofile' */
				fallthrough

			case "update":

				/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
				if tokeninfo == nil {

					/* This  function return error codes to html */
					FullData["errorresults"] = functions.ReturnError("Warning", "You must be Signed In")

					/* This is the template execution for 'index' */
					functions.ExecuteTemplate(w, "index.html", FullData)

					break

				} else {

					/* Execute all calculations and returns err if validation fails */
					var err = calculations.Calculate(FullData, r)

					/* This conditional does not allow transactions to be ecorded on the
					database when there is a validation error or profile change*/
					if err == false && r.PostFormValue("submitselect") != "vmprofile" {

						/* This function save the transaction into vdicalc.transactions */
						mysql.SaveTransaction(db, tokeninfo.UserId, functions.GetIP(r), FullData)

					}
				}

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", FullData)

			}

		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}

	}

	// db.Close()

}
