package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/New-Era/servers/gateway/handlers"
	"github.com/New-Era/servers/gateway/models/devices"
	"github.com/New-Era/servers/gateway/sessions"
	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
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

	sessionKey := os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.Fatal("please set session key")
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
	// dsn := fmt.Sprintf("root:%s@tcp(mysql:3306)/mysql",
	// 	os.Getenv("MYSQL_ROOT_PASSWORD"))
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	fmt.Printf("error opening database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer db.Close()

	sessStore := sessions.NewRedisStore(rClient, time.Duration(600)*time.Second)
	deviceStore := devices.NewMongoStore(mongoSess)
	conn := handlers.NewConnections()
	hc := handlers.NewHandlerContext(sessionKey, sessStore, deviceStore, conn)

	// addresses of websocket microservice instances
	wc := strings.Split(os.Getenv("WCADDRS"), ",")

	// proxy for websocket microservice
	wcProxy := &httputil.ReverseProxy{Director: CustomDirectorRR(wc, hc)}

	mux := http.NewServeMux()
	mux.HandleFunc("/time", handlers.TimeHandler)
	mux.HandleFunc("/device", hc.DevicesHandler)
	mux.Handle("/ws", wcProxy)
	wrappedMux := handlers.NewCORS(mux)
	fmt.Printf("server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedMux))
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
		_, _ = sessions.GetState(r, hc.SigningKey, hc.SessStore, &tmp)
		if tmp.Device.ID != "" { // set if user exists
			j, err := json.Marshal(tmp.Device)
			if err != nil {
				fmt.Println(err)
				return
			}
			r.Header.Set("X-Device", string(j))
		}
		r.Header.Add("X-Forwarded-Host", r.Host)
		if strings.HasPrefix(dest.String(), "wc") {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			r.URL.Scheme = "https"
		} else {
			r.URL.Scheme = "http"
		}
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
		_, _ = sessions.GetState(r, hc.SigningKey, hc.SessStore, &tmp)
		if tmp.Device.ID != "" { // set if user exists
			j, err := json.Marshal(tmp.Device)
			if err != nil {
				fmt.Println(err)
				return
			}
			r.Header.Set("X-Device", string(j))
		}
		r.Header.Add("X-Forwarded-Host", r.Host)
		if strings.HasPrefix(target.String(), "wc") {
			fmt.Println(r)
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			r.URL.Scheme = "https"
			r.Header.Set("X-Connection", "Upgrade")
			r.Header.Set("X-Upgrade", "Websocket")
			r.Header.Set("X-Sec-Websocket-Key", r.Header.Get("Sec-Websocket-Key"))
		} else {
			r.URL.Scheme = "http"
		}
		r.URL.Host = target.String()
		r.Host = target.String()
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
