package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	c "vdicalc/config"
	f "vdicalc/functions"
	host "vdicalc/host"
	"vdicalc/mysql"
	storage "vdicalc/storage"
	v "vdicalc/virtualization"

	"github.com/spf13/viper"
)

var tlp *template.Template
var configuration c.Configurations
var db *sql.DB

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

		/* This is the template execution for 'index' */
		err := tlp.ExecuteTemplate(w, "index.html", fullData)
		if err != nil {
			panic(err)
		}

	case "POST":

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
			fullData["vmclonesizerefreshrateselected"] = key.Clonesizerefreshrateselected

			/* Fallthrough ensures 'update' is always executed after 'vmprofile' */
			fallthrough

		case "update":

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

			/* This is the template execution for 'update' */
			err := tlp.ExecuteTemplate(w, "index.html", fullData)
			if err != nil {
				panic(err)
			}

			/* This conditional does not allow profile changes to be recorded on the database */
			if r.PostFormValue("submitselect") != "vmprofile" {

				// If the optional DB_TCP_HOST environment variable is set, it contains
				// the IP address and port number of a TCP connection pool to be created,
				// such as "127.0.0.1:3306". If DB_TCP_HOST is not set, a Unix socket
				// connection pool will be created instead.
				if os.Getenv("DB_TCP_HOST") != "" {

					db, err = mysql.InitTCPConnectionPool()
					if err != nil {
						log.Fatalf("initTCPConnectionPool: unable to connect: %v", err)
					}
				} else {
					db, err = mysql.InitSocketConnectionPool()
					if err != nil {
						log.Fatalf("initSocketConnectionPool: unable to connect: %v", err)
					}
				}

				/* Build MySQL statement  */
				sqlInsert, _ := mysql.SQLBuilder(f.GetIP(r), fullData["hostresultscount"], fullData["hostresultsclockused"], fullData["hostresultsmemory"], fullData["hostresultsvmcount"], fullData["storageresultscapacity"], fullData["storageresultsdatastorecount"], fullData["storageresultsdatastoresize"], fullData["storagedatastorefroentendiops"], fullData["storagedatastorebackendiops"], fullData["storageresultsfrontendiops"], fullData["storageresultsbackendiops"])

				/* This function execues the SQL estatement on Google SQL Run database */
				mysql.Insert(db, sqlInsert)
			}

		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}
