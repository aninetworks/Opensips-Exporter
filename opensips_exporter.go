package main

/* Changelog
 * S6760 - 05/05/2020 - Mario - Initial version.
 * This application get a JSON file from the OpenSIPS server and translate to a metric in order to be used on Prometheus
 * Then it will start hosting a webserver one or two instances of it according the config file when the URL is NOT empty
 * Prometheus must be configurated to point to the webserver and the ports must be open on the firewall.
 * Grafana must use Prometheus as datasource.
 * 2020-10-02 - Mario - Changed to serve multiple servers and both OpenSips versions (2.4 & 3.1)
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"CustomOpenSipsExporter/bin/osipsclasses"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

// MyType used to store the info of the JSON
type MyType struct {
	field  string
	value  string
	desc   string
	metric string
}

// ConfigType used to get the configuration from the yaml file
type ConfigType struct {
	PIDFile struct {
		Path string `yaml:"path"`
		Name string `yaml:"name"`
	} `yaml:"PIDFile"`
	Server []struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		URL     string `yaml:"URL"`
		Body    string `yaml:"body"`
		Version string `yaml:"version"`
	} `yaml:"Server"`
}

// OpensipResult type
type OpensipResult struct {
	Jsonrpc string           `json:"jsonrpc"`
	Result  map[string]int64 `json:"result"`
	ID      int64            `json:"id"`
}

// Opensiptrunkcalls struct
type Opensiptrunkcalls struct {
	Jsonrpc string   `json:"jsonrpc"`
	Result  []Result `json:"result"`
	ID      int64    `json:"id"`
}

// Result struct
type Result struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

var cfg ConfigType

const regexMetric = `^[a-zA-Z_:][a-zA-Z0-9_:]*$`

// Location for the config file requested by Joseph
const configFilePath = "/opt/ani_opensips_exporter/etc/"

func getContentUsingGJson(url string, body string, version string) []*MyType {
	tmpDataArray := []*MyType{}

	if version == "2" {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		if !gjson.Valid(string(jsonDataFromHTTP)) {
			panic("invalid JSON format")
		}

		m, ok := gjson.Parse(string(jsonDataFromHTTP)).Value().(map[string]interface{})
		if !ok {
			fmt.Println("Error parsing JSON")
		}

		for k, v := range m {
			switch vv := v.(type) {
			case string:
				stat := strings.Split(k, ":")
				retVal := osipsclasses.Metric{}
				tmpData := new(MyType)
				if strings.ToLower(stat[0]) != "snmpstats" {
					if len(stat) >= 2 {
						switch strings.ToLower(stat[0]) {
						case "core":
							retVal = osipsclasses.CoreDescription(stat[1])
						case "load":
							retVal = osipsclasses.LoadDescription(stat[1])
						case "net":
							retVal = osipsclasses.NetDescription(stat[1])
						case "pkmem":
							retVal = osipsclasses.PkmemDescription(stat[1])
						case "shmem":
							retVal = osipsclasses.ShmemDescription(stat[1])
						default:
							retVal.Desc = "Unrecognized metric."
							retVal.Value = "gauge"
						}
					}
					tmpData.field = stat[1]
					tmpData.value = vv
					tmpData.desc = retVal.Desc
					tmpData.metric = retVal.Value
					tmpDataArray = append(tmpDataArray, tmpData)
				}
			case []interface{}:
				for _, u := range vv {
					if reflect.ValueOf(u).Kind() == reflect.Map {
						d := reflect.ValueOf(u)
						tmpData := new(MyType)
						for _, k := range d.MapKeys() {
							if tmpData.field != "" && tmpData.value != "" {
								tmpData.field = ""
								tmpData.value = ""
								tmpData.desc = ""
								tmpData.metric = ""
							}
							if reflect.ValueOf(d.MapIndex(k).Interface()).Kind() == reflect.String {
								tmpData.field = d.MapIndex(k).Interface().(string)
							} else if reflect.ValueOf(d.MapIndex(k).Interface()).Kind() == reflect.Map {
								g := reflect.ValueOf(d.MapIndex(k).Interface())
								for _, z := range g.MapKeys() {
									tmpData.value = g.MapIndex(z).Interface().(string)
								}
							}
							if tmpData.value != "" && tmpData.field != "" {
								tmpData.desc = "Trunkgroup Name"
								tmpData.metric = "gauge"
								tmpDataArray = append(tmpDataArray, tmpData)
							}
						}
					}
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}
	} else {
		var jsonStr = []byte(body)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		jsonDataFromHTTP, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		if !gjson.Valid(string(jsonDataFromHTTP)) {
			panic("invalid JSON format")
		}

		if strings.Contains(body, "get_statistics") {
			var r OpensipResult
			errm := json.Unmarshal(jsonDataFromHTTP, &r)
			if err != nil {
				panic(errm)
			}

			for k, v := range r.Result {
				stat := strings.Split(k, ":")
				retVal := osipsclasses.Metric{}
				tmpData := new(MyType)
				if strings.ToLower(stat[0]) != "snmpstats" {
					if len(stat) >= 2 {
						switch strings.ToLower(stat[0]) {
						case "core":
							retVal = osipsclasses.CoreDescription(stat[1])
						case "load":
							retVal = osipsclasses.LoadDescription(stat[1])
						case "net":
							retVal = osipsclasses.NetDescription(stat[1])
						case "pkmem":
							retVal = osipsclasses.PkmemDescription(stat[1])
						case "shmem":
							retVal = osipsclasses.ShmemDescription(stat[1])
						default:
							retVal.Desc = "Unrecognized metric."
							retVal.Value = "gauge"
						}
					}
					tmpData.field = stat[1]
					tmpData.value = strconv.FormatInt(v, 10)
					tmpData.desc = retVal.Desc
					tmpData.metric = retVal.Value
					tmpDataArray = append(tmpDataArray, tmpData)
				}
			}
		} else {
			var r Opensiptrunkcalls
			errm := json.Unmarshal(jsonDataFromHTTP, &r)
			if err != nil {
				panic(errm)
			}

			for _, v := range r.Result {
				tmpData := new(MyType)
				tmpData.field = v.Value
				tmpData.value = strconv.FormatInt(v.Count, 10)
				tmpData.desc = "Trunkgroup Name"
				tmpData.metric = "gauge"
				tmpDataArray = append(tmpDataArray, tmpData)
			}
		}
	}
	return tmpDataArray
}

func main() {
	// Configure logger to write to the syslog.
	// Comment these section for Windows since syslog does not exists on Windows
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "opensips_exporter")
	if e == nil {
		log.SetOutput(logwriter)
	}

	//Get filename
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Unable to get own process name.")
		panic("Unable to get own process name.")
	}

	//Get the file name to look for a config file with the same name
	fileName := strings.Split(filepath.Base(file), ".")[0]

	//Open the config file
	f, err := os.Open(configFilePath + fileName + ".yml")
	if err != nil {
		log.Println("err")
		panic(err)
	}
	defer f.Close()

	//Load the config file
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Println("err")
		panic(err)
	}

	// //Write the PID file (LINUX only)
	err = osipsclasses.WritePidFile(cfg.PIDFile.Path, cfg.PIDFile.Name)
	if err != nil {
		log.Println("err")
		panic(err)
	}

	// //Delete old PID file if exists (LINUX only)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		deletePIDFile(cfg.PIDFile.Path + cfg.PIDFile.Name)
		os.Exit(0)
	}()	

	finish := make(chan bool)
	for _, s := range cfg.Server {		
		//Start the process if the URL is not empty
		if s.URL != "" {
			go callSubroutine(s.URL, s.Host, s.Port, s.Body, s.Version)
			time.Sleep(1000 * time.Millisecond)
		}
	}
	<-finish
}

func callSubroutine(url string, host string, port string, body string, version string) {
	ServerMux := http.NewServeMux()
	ServerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var validMetric = regexp.MustCompile(regexMetric)
		for _, y := range getContentUsingGJson(url, body, version) {
			y.field = strings.ReplaceAll(y.field, "-", "_")
			y.field = strings.ReplaceAll(y.field, " ", "_")
			y.field = strings.ReplaceAll(y.field, ".", "")
			if !validMetric.MatchString(y.field) {
				y.field = "id_" + y.field
			}
			fmt.Fprintf(w, "# HELP %s %s\n", y.field, y.desc)
			fmt.Fprintf(w, "# TYPE %s %s\n", y.field, y.metric)
			fmt.Fprintf(w, "%s %s\n", y.field, y.value)
		}
	})
	log.Printf("Started Custom OpenSIPS exporter, listening on %v", (host + ":" + port))
	log.Fatal(http.ListenAndServe((host + ":" + port), ServerMux))
}

func deletePIDFile(fileName string) {
	os.Remove(fileName)
}

