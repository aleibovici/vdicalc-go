package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"vdicalc/auth"
	"vdicalc/config"
	c "vdicalc/config"
	f "vdicalc/functions"
	host "vdicalc/host"
	"vdicalc/mysql"
	storage "vdicalc/storage"
	"vdicalc/validation"
	v "vdicalc/virtualization"

	"github.com/spf13/viper"
	"google.golang.org/api/oauth2/v2"
)

var tlp *template.Template
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

	var err error

	tlp, err = template.ParseGlob("./templates/*")
	if err != nil {
		panic(err)
	}

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

	/* This function call up a function (functions/Dataload) to dynamicaly populate the dataset for the HTML files */
	fullData := f.DataLoad(configuration)

	switch r.Method {
	case "GET":

		/* Determine the backend service address to be passsed to index.html js for authentication.
		AUTH_ADDRESS is an environment variable and must be defined at the container execution level */
		fullData["authaddress"] = os.Getenv("AUTH_ADDRESS")

		/* This is the template execution for 'index' */
		err := tlp.ExecuteTemplate(w, "index.html", fullData)
		if err != nil {
			panic(err)
		}

	case "POST":

		/* Determine the URI path to de taken */
		switch r.URL.Path {

		case "/tokensignin":

			/* The index.html signin JS triggers a http POST on /tokensignin and the user's token is associated to tokeninfo for verification*/
			if tokeninfo == nil {

				tokeninfo, _ = auth.VerifyIDToken(strings.TrimPrefix(r.URL.RawQuery, "id_token="))

				/* Inititalize DB connection */
				db = mysql.DBInit()

				/* Test if user exist in vdicalc.users table and if not add user to the database*/
				if mysql.QueryUser(db, tokeninfo.UserId) == false {

					/* This function executes the SQL estatement on Google SQL Run database */
					mysql.CreateUser(db, tokeninfo.UserId, tokeninfo.Email)

				}
			}

		case "/tokensignoff":

			/* The index.html signoff JS triggers a http POST on /tokensignoff and the user's token id is removed*/
			tokeninfo = nil

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
				fullData[newKey] = values[0]

			}

			/* This function uses a hidden field 'submitselect' in each HTML template to detect the actions triggered by users.
			HTML action must include 'document.getElementById('submitselect').value='about';this.form.submit()' */
			switch r.PostFormValue("submitselect") {

			case "about":

				/* This is the template execution for 'about' */
				err := tlp.ExecuteTemplate(w, "about.html", "")
				if err != nil {
					panic(err)
				}

			case "back", "guide":

				/* This is the template execution for 'back' */
				err := tlp.ExecuteTemplate(w, "index.html", fullData)
				if err != nil {
					panic(err)
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

				fullData["vmvcpucountselected"] = key.Vcpucountselected
				fullData["vmvcpumhzselected"] = key.Vcpumhz
				fullData["vmpercorecountselected"] = key.Vmpercorecountselected
				fullData["vmmemorysizeselected"] = key.Memorysize
				fullData["vmdisplaycountselected"] = key.Displaycountselected
				fullData["vmdisplayresolutionselected"] = key.Displayresolutionselected
				fullData["vmvideoramselected"] = key.Videoramselected
				fullData["vmdisksizeselected"] = key.Disksizeselected
				fullData["vmiopscountselected"] = key.Iopscountselected
				fullData["vmiopsreadratioselected"] = key.Iopsreadratioselected
				fullData["vmclonesizerefreshrateselected"] = key.Clonesizerefreshrateselected

				/* Fallthrough ensures 'update' is always executed after 'vmprofile' */
				fallthrough

			case "update":

				/* If not valid tokenID is present a error is raised requesting user to signin, else execute calculations */
				if tokeninfo == nil {

					var errorList config.ErrorResultsConfiguration
					error := config.ErrorConfiguration{Code: "Warning: ", Description: "You must be Signed In"}
					errorList.Error = append(errorList.Error, error)
					fullData["errorresults"] = errorList

					/* This is the template execution for 'update' */
					err := tlp.ExecuteTemplate(w, "index.html", fullData)
					if err != nil {
						panic(err)
					}

					break

				} else {

					/* This is the default template execution mode with results calculation */
					fullData["hostresultscount"] = host.GetHostCount(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"))
					fullData["hostresultsclockused"] = host.GetHostClockUsed(r.FormValue("vmvcpucount"), r.FormValue("vmvcpumhz"), r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"))
					fullData["hostresultsmemory"] = host.GetHostMemory(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("hostcoresoverhead"), r.FormValue("vmpercorecount"), r.FormValue("vmmemorysize"), r.FormValue("hostmemoryoverhead"), r.FormValue("vmdisplaycount"), r.FormValue("vmdisplayresolution"), r.FormValue("vmvcpucount"), r.FormValue("vmvideoram"))
					fullData["hostresultsvmcount"] = host.GetHostVMCount(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"))
					fullData["storageresultscapacity"] = storage.GetStorageCapacity(r.FormValue("vmcount"), r.FormValue("vmdisksize"), r.FormValue("storagecapacityoverhead"), r.FormValue("storagededuperatio"), r.FormValue("vmdisplaycount"), r.FormValue("vmdisplayresolution"), r.FormValue("vmvideoram"), r.FormValue("vmmemorysize"), r.FormValue("vmclonesizerefreshrate"))
					fullData["storageresultsdatastorecount"] = storage.GetStorageDatastoreCount(r.FormValue("vmcount"), r.FormValue("storagedatastorevmcount"))
					fullData["storageresultsdatastoresize"] = storage.GetStorageDatastoreSize(r.FormValue("vmcount"), r.FormValue("storagedatastorevmcount"), r.FormValue("vmdisksize"), r.FormValue("storagecapacityoverhead"), r.FormValue("storagededuperatio"), r.FormValue("vmdisplaycount"), r.FormValue("vmdisplayresolution"), r.FormValue("vmvideoram"), r.FormValue("vmmemorysize"), r.FormValue("vmclonesizerefreshrate"))
					fullData["storagedatastorefroentendiops"], fullData["storagedatastorebackendiops"], fullData["storageresultsfrontendiops"], fullData["storageresultsbackendiops"] = storage.GetStorageDatastoreIops(r.FormValue("vmiopscount"), r.FormValue("vmiopsreadratio"), r.FormValue("storagedatastorevmcount"), r.FormValue("storageraidtype"), r.FormValue("vmcount"), r.FormValue("storagedatastorevmcount"))
					fullData["virtualizationresultsclustercount"] = v.GetClusterSize(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"), r.FormValue("virtualizationclusterhostsize"))
					fullData["virtualizationresultsmanagementservercount"] = v.GetManagementServerCount(r.FormValue("vmcount"), r.FormValue("virtualizationmanagementservertvmcount"))
					fullData["errorresults"] = validation.ValidateResults(fullData)
				}

				/* This is the template execution for 'update' */
				err := tlp.ExecuteTemplate(w, "index.html", fullData)
				if err != nil {
					panic(err)
				}

				/* This conditional does not allow profile changes to be recorded on the database */
				if r.PostFormValue("submitselect") != "vmprofile" {

					/* Build MySQL statement  */
					sqlInsert, _ := mysql.SQLBuilderInsert("vdicalc", map[string]interface{}{
						"datetime":                      time.Now(),
						"guserid":                       tokeninfo.UserId,
						"ip":                            f.GetIP(r),
						"hostresultscount":              fmt.Sprint(fullData["hostresultscount"]),
						"hostresultsclockused":          fmt.Sprint(fullData["hostresultsclockused"]),
						"hostresultsmemory":             fmt.Sprint(fullData["hostresultsmemory"]),
						"hostresultsvmcount":            fmt.Sprint(fullData["hostresultsvmcount"]),
						"storageresultscapacity":        fmt.Sprint(fullData["storageresultscapacity"]),
						"storageresultsdatastorecount":  fmt.Sprint(fullData["storageresultsdatastorecount"]),
						"storageresultsdatastoresize":   fmt.Sprint(fullData["storageresultsdatastoresize"]),
						"storagedatastorefroentendiops": fmt.Sprint(fullData["storagedatastorefroentendiops"]),
						"storagedatastorebackendiops":   fmt.Sprint(fullData["storagedatastorebackendiops"]),
						"storageresultsfrontendiops":    fmt.Sprint(fullData["storageresultsfrontendiops"]),
						"storageresultsbackendiops":     fmt.Sprint(fullData["storageresultsbackendiops"]),
					})

					/* This function execues the SQL estatement on Google SQL Run database */
					mysql.Insert(db, sqlInsert)
				}

			}

		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}

	}

}
