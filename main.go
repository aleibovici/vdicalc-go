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
	c "vdicalc/config"
	"vdicalc/functions"
	"vdicalc/mysql"

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

		}

	case "POST":

		/* Determine the URI path to de taken */
		switch r.URL.Path {

		case "/tokensignoff":

			/* The index.html signoff JS triggers a http POST on /tokensignoff and the user's token id is removed*/
			tokeninfo = nil

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

						/* This function build a SQL statement for inserting calculation data into vdicalc.vdicalc  */
						sqlInsert, _ := mysql.SQLBuilderInsert("saves", map[string]interface{}{
							"datetime":                                       time.Now(),
							"guserid":                                        tokeninfo.UserId,
							"savename":                                       strings.ToUpper((time.Now().Format("01-02-2006 15:04:05")) + " " + (r.PostFormValue("savename"))),
							"vmcountselected":                                fmt.Sprint(FullData["vmcountselected"]),
							"vmvcpucountselected":                            fmt.Sprint(FullData["vmvcpucountselected"]),
							"vmvcpumhzselected":                              fmt.Sprint(FullData["vmvcpumhzselected"]),
							"vmpercorecountselected":                         fmt.Sprint(FullData["vmpercorecountselected"]),
							"vmdisplaycountselected":                         fmt.Sprint(FullData["vmdisplaycountselected"]),
							"vmdisplayresolutionselected":                    fmt.Sprint(FullData["vmdisplayresolutionselected"]),
							"vmmemorysizeselected":                           fmt.Sprint(FullData["vmmemorysizeselected"]),
							"vmvideoramselected":                             fmt.Sprint(FullData["vmvideoramselected"]),
							"vmdisksizeselected":                             fmt.Sprint(FullData["vmdisksizeselected"]),
							"vmiopscountselected":                            fmt.Sprint(FullData["vmiopscountselected"]),
							"vmiopsreadratioselected":                        fmt.Sprint(FullData["vmiopsreadratioselected"]),
							"vmclonesizerefreshrateselected":                 fmt.Sprint(FullData["vmclonesizerefreshrateselected"]),
							"hostsocketcountselected":                        fmt.Sprint(FullData["hostsocketcountselected"]),
							"hostsocketcorescountselected":                   fmt.Sprint(FullData["hostsocketcorescountselected"]),
							"hostmemoryoverheadselected":                     fmt.Sprint(FullData["hostmemoryoverheadselected"]),
							"hostcoresoverheadselected":                      fmt.Sprint(FullData["hostcoresoverheadselected"]),
							"storagecapacityoverheadselected":                fmt.Sprint(FullData["storagecapacityoverheadselected"]),
							"storagedatastorevmcountselected":                fmt.Sprint(FullData["storagedatastorevmcountselected"]),
							"storagededuperatioselected":                     fmt.Sprint(FullData["storagededuperatioselected"]),
							"storageraidtypeselected":                        fmt.Sprint(FullData["storageraidtypeselected"]),
							"virtualizationclusterhostsizeselected":          fmt.Sprint(FullData["virtualizationclusterhostsizeselected"]),
							"virtualizationmanagementservertvmcountselected": fmt.Sprint(FullData["virtualizationmanagementservertvmcountselected"]),
							"virtualizationclusterhosthaselected":            fmt.Sprint(FullData["virtualizationclusterhosthaselected"]),
						})

						/* This function execues the SQL estatement on Google SQL Run database */
						mysql.Insert(db, sqlInsert)

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

					calculations.Calculate(FullData, r)
				}

				/* This is the template execution for 'index' */
				functions.ExecuteTemplate(w, "index.html", FullData)

				/* This conditional does not allow profile changes to be recorded on the database */
				if r.PostFormValue("submitselect") != "vmprofile" {

					/* This function build a SQL statement for inserting calculation data into vdicalc.vdicalc  */
					sqlInsert, _ := mysql.SQLBuilderInsert("transactions", map[string]interface{}{
						"datetime":                      time.Now(),
						"guserid":                       tokeninfo.UserId,
						"ip":                            functions.GetIP(r),
						"hostresultscount":              fmt.Sprint(FullData["hostresultscount"]),
						"hostresultsclockused":          fmt.Sprint(FullData["hostresultsclockused"]),
						"hostresultsmemory":             fmt.Sprint(FullData["hostresultsmemory"]),
						"hostresultsvmcount":            fmt.Sprint(FullData["hostresultsvmcount"]),
						"storageresultscapacity":        fmt.Sprint(FullData["storageresultscapacity"]),
						"storageresultsdatastorecount":  fmt.Sprint(FullData["storageresultsdatastorecount"]),
						"storageresultsdatastoresize":   fmt.Sprint(FullData["storageresultsdatastoresize"]),
						"storagedatastorefroentendiops": fmt.Sprint(FullData["storagedatastorefroentendiops"]),
						"storagedatastorebackendiops":   fmt.Sprint(FullData["storagedatastorebackendiops"]),
						"storageresultsfrontendiops":    fmt.Sprint(FullData["storageresultsfrontendiops"]),
						"storageresultsbackendiops":     fmt.Sprint(FullData["storageresultsbackendiops"]),
					})

					/* This function execues the SQL estatement on Google SQL Run database */
					mysql.Insert(db, sqlInsert)
				}

			}

		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}

	}

	db.Close()

}
