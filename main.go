/*
	Services: pingservices
	Author: Kjell Osse Almgren, Tetracon AB
	Date: 2017-08-16
	Description: Service to check services availability in project Lagerlöf
	Architecture:
	win32: GOOS=windows GOARCH=386 go build -v
	win64: GOOS=windows GOARCH=amd64 go build -v
	arm64: GOOS=linux GOARCH=arm64 go build -v
	amd64: GOOS=linux GOARCH=amd64 go build -v
	arm: GOOS=linux GOARCH=arm go build -v
	exprimental: GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags '-w -s' -a -installsuffix cgo -o pingservices
	expriemntal: CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags pingservices -ldflags '-w'
	exprimental: GOOS=linux GOARCH=arm64 go build -a --ldflags 'extldflags "-static"' -tags pingservices -installsuffix pingservices .
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pingservices/version"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Banner text
const (
	// TETRACON banner
	TETRACON = `
_________    __
|__    __|   | |
   |  |  ___ | |_   ____  ___   ___ ___  _ __ 
   |  | / _ \|  _| /  __|/ _ \ / __/ _ \| '_ \
   |  | \ __/| |_  | |  | (_| | (_| (_) | | | | 
   |__| \___| \__| |_|   \__,_|\___\___/|_| |_| 
version: %s
`
)

//
var (
	srv  bool
	vrsn bool
)

var (
	arg string
)

// template
var tpl *template.Template

//
//	JSON struct for configuration file
//
type service struct {
	Target      string `json:"target"`
	Environment string `json:"environment"`
	Urlstring   string `json:"urlstring"`
	Contact     string `json:"contact"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

//
type Ping struct {
	Target      string
	Environment string
	Urlstring   string
	Contact     string
	Email       string
	Phone       string
	Ping        bool
	Errstring   string
	Httpcode    int
}

//
type MyPinglists struct {
	Hostname string
	Pings    []Ping
}

// init function
func init() {

	// instanciate a new logger
	var log = logrus.New()
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&srv, "server", true, "run in server mode")
	flag.BoolVar(&srv, "s", true, "run in server mode (shorthand)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(TETRACON, version.PingVersion()))
		flag.PrintDefaults()
	}

	flag.Parse()
	if vrsn {

		fmt.Printf("flag version %s\n", version.PingVersion())
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		arg = flag.Args()[0]
	}

	if arg == "help" {
		usageAndExit("", 0)
	}

	if arg == "version" {
		fmt.Printf("flag version %s\n", version.PingVersion())
		os.Exit(0)
	}

	//
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	//log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter) // default

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	log.Level = logrus.DebugLevel
}

func main() {

	port := 9000
	//
	//	Read json configuration file
	//
	// parse the arg
	//arg := flag.Args()[0]
	//
	// check both possible arguments
	if flag.NArg() < 1 {
		showStartup(port)
		color.Unset()
		router := mux.NewRouter()
		router.HandleFunc("/health-check", HealthCheckHandler).Methods("GET")
		router.HandleFunc("/pingqa", PingHandler).Methods("GET")
		router.HandleFunc("/pingprod", PingHandler).Methods("GET")

		router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
		router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
		router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))
		//http.Handle("/", router)
		//
		err := http.ListenAndServe(":"+strconv.Itoa(port), router)
		if err != nil {
			fmt.Printf("ListenAndServer Error: %s", err.Error())
			logrus.Fatal(err)
		}
	}
}

//
//	Get hostname of running server
//
func GetHostname() string {

	hostname, err1 := os.Hostname()
	if err1 == nil {
		//log.Printf("Hostname: %s", hostname)
		//fmt.Println("Hostname: ", hostname)
	}
	return hostname
}

//
//	just for health-check, can be removed
//
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
	io.WriteString(w, `{"status": http.StatusOK}`)
	fmt.Printf("Http-Status %d received\r\n", http.StatusOK)
}

//
//	Build list array to ping execute (Pings)
//
func (pings *MyPinglists) AddItem(ping Ping) []Ping {
	pings.Hostname = ""
	pings.Pings = append(pings.Pings, ping)
	return pings.Pings
}

//
// handler for ping requests
func PingHandler(w http.ResponseWriter, r *http.Request) {

	filePath := ""
	param := strings.Split(r.URL.Path, "/")
	//fmt.Printf("URL=%s /r/n", r.URL.Path)
	// r.URL.Path = /pingqs or /pingprod
	fmt.Printf("Len=%d", len(param))
	par := ""
	if len(param) == 2 {
		par = param[1]
	}
	//
	if strings.Contains(par, "pingprod") {
		filePath = "./services-prod.json"
	} else if strings.Contains(par, "pingqa") {
		filePath = "./services-qa.json"
	}

	//filePath := "./services-prod.json"
	file, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		fmt.Printf("Error reading configuration file %s\r\n", filePath)
		fmt.Printf("File error: %v\r\n", err1)
		os.Exit(1)
	}
	var services []service
	err2 := json.Unmarshal(file, &services)
	if err2 != nil {
		fmt.Printf("JSON marshal Error: %s\r\n", err2)
		fmt.Printf("Check %s for JSON typing error\r\n", filePath)
		os.Exit(1)
	}
	pings := []Ping{} // Initialize
	i := MyPinglists{GetHostname(), pings}
	//
	fmt.Printf("Hostname Kjell: %s", i.Hostname)
	for key := range services {

		fmt.Printf("Processing target (")
		color.Set(color.FgHiWhite)
		fmt.Printf("%s", services[key].Target)
		color.Unset()
		fmt.Printf(") url ")
		color.Set(color.FgHiWhite)
		fmt.Printf("%s - ", services[key].Urlstring)
		color.Unset()
		httpcode, err := PingExec(services[key].Target, services[key].Urlstring)
		if err == nil {
			if httpcode == 200 {
				color.Set(color.FgHiGreen)
				fmt.Printf("%d", httpcode)
				color.Unset()
				fmt.Printf(" OK\r\n")
				ping := Ping{services[key].Target, services[key].Environment, services[key].Urlstring, services[key].Contact, services[key].Email, services[key].Phone, true, "OK", httpcode}
				i.AddItem(ping)
			}
		}
		if httpcode == 418 {
			//color.Set(color.FgYellow)
			//fmt.Printf(" %s", ping.Errstring)
			//fmt.Println("")
			//color.Unset()
			ping := Ping{services[key].Target, services[key].Environment, services[key].Urlstring, services[key].Contact, services[key].Email, services[key].Phone, true, "OK", httpcode}
			//items := []Test{}
			//tests = MyTests{items}
			//i := MyInventories{inventories}
			i.AddItem(ping)
		}
		if httpcode == 401 {
			fmt.Printf("service ")
			color.Set(color.FgHiRed)
			fmt.Printf(" %s", services[key].Target)
			color.Unset()
			fmt.Printf(" received 401\r\n")
			ping := Ping{services[key].Target, services[key].Environment, services[key].Urlstring, services[key].Contact, services[key].Email, services[key].Phone, false, "Unauthorized access", httpcode}
			i.AddItem(ping)
		}
		if httpcode >= 500 {
			fmt.Printf("service ")
			color.Set(color.FgHiRed)
			fmt.Printf("%s", services[key].Target)
			color.Unset()
			fmt.Printf(" unavailable\r\n")
			color.Unset()
			ping := Ping{services[key].Target, services[key].Environment, services[key].Urlstring, services[key].Contact, services[key].Email, services[key].Phone, false, "Unavailable", httpcode}
			i.AddItem(ping)
		}
	}
	i.Hostname = GetHostname()
	//
	wpage := ""
	if strings.Contains(par, "pingprod") {
		wpage = "index-prod.html"
	} else if strings.Contains(par, "pingqa") {
		wpage = "index-qa.html"
	}
	err := tpl.ExecuteTemplate(w, wpage, i)
	if err != nil {
		log.Fatal(err)
	}
}

// usageAndExit
func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}

// showStartup
func showStartup(port int) {
	color.Set(color.FgHiGreen)
	fmt.Printf("Lagerlöf availability services is started on server: ")
	color.Set(color.FgHiWhite)
	fmt.Printf("%s", GetHostname())
	color.Set(color.FgHiGreen)
	fmt.Printf(" is listen on port ")
	color.Set(color.FgHiWhite)
	fmt.Printf("%d", port)
	color.Set(color.FgHiGreen)
	fmt.Println(" tls")
	color.Unset()
}
