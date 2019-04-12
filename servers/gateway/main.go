package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/New-Era/servers/gateway/handlers"
	"github.com/New-Era/servers/gateway/models/devices"
	"github.com/New-Era/servers/gateway/sessions"
	mgo "gopkg.in/mgo.v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

// main entry point for the server
func main() {
	// Connection to HTTPS
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	sessionKey = os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.fatal("please set session key")
	}

	tlscert := os.Getenv("TLSCERT")
	tlskey := os.Getenv("TLSKEY")
	if len(tlskey) == 0 && len(tlscert) == 0 {
		log.Fatal("please set TLSKEY and TLSCERT")
	}

	// REDIS DB ADDRESS
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = ":6397"
	}

	rClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	// MONGO DB CONNECTION
	// get the address of the MongoDB server from an environment variable
	mongoAddr := os.Getenv("MONGO_ADDR")
	//default to "localhost"
	if len(mongoAddr) == 0 {
		mongoAddr = "localhost"
	}
	// Dialing MongoDB server
	mongoSess, err := mgo.Dial(mongoAddr)
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	// TODO: construct a new MongoStore, provide mongoSess as well as a
	// database and collection name to use (device maybe?)

	// MYSQL DB CONNECTION
	// Construct MySql serve
	dsn := fmt.Sprintf("root:%s@tcp(mysql:3306)/mysql",
		os.Getenv("MYSQL_ROOT_PASSWORD"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	sessStore = sessions.NewRedisStore(rClient, time.Duration(600)*time.Second)
	deviceStore := devices.NewMongoStore(mongoSess)
	conn := handlers.NewConnections()
	handlerCtx := handlers.NewHandlerContext(sessionKey, sessStore, deviceStore, conn)

	// addresses of websocket microservice instances
	wc := strings.Split(os.Getenv("WCADDRS"), ",")

	// proxy for websocket microservice
	wcProxy := &httputil.ReverseProxy{Director: CustomDirectorRR(wc, &hc)}

	mux := http.NewServeMux()
	mux.HandleFunc("/time", handlers.TimeHandler)
	mux.HandleFunc("/device", handlerCtx.DevicesHandler)

	fmt.Printf("server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, mux))
}

// Director handles the transport of requests to proper endpoints
type Director func(r *http.Request)

// CustomDirectorRR directs requests for services that have
// multiple servers via Round Robin technique
func CustomDirectorRR(targets []string, hc *handlers.HandlerContext) Director {
	if len(targets) == 1 {
		dest, _ := url.Parse(targets[0])
		return CustomDirector(dest, hc)
	}
	var i int32
	i = 0
	url, _ := url.Parse(targets[int(i)%len(targets)])
	atomic.AddInt32(&i, 1)
	dest := url
	return func(r *http.Request) {
		r.Header.Del("X-Device") // remove any previous user
		tmp := handlers.SessionState{}
		_, _ = s.GetState(r, hc.Key, hc.Session, &tmp)
		if tmp.Device.ID != 0 { // set if user exists
			j, err := json.Marshal(tmp.Device)
			if err != nil {
				fmt.Println(err)
				return
			}
			r.Header.Set("X-Device", string(j))
		}
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.URL.Scheme = "http"
		r.URL.Host = dest.String()
		r.Host = dest.String()
	}
}

// CustomDirector directs requests to a specified server and modifies the request
// before being passed along
func CustomDirector(target *url.URL, hc *handlers.HandlerContext) Director {
	return func(r *http.Request) {
		r.Header.Del("X-Device") // remove any previous user
		tmp := handlers.SessionState{}
		_, _ = s.GetState(r, hc.Key, hc.Session, &tmp)
		if tmp.Device.ID != 0 { // set if user exists
			j, err := json.Marshal(tmp.User)
			if err != nil {
				fmt.Println(err)
				return
			}
			r.Header.Set("X-Device", string(j))
		}
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.URL.Scheme = "http"
		r.URL.Host = target.String()
		r.Host = target.String()
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
