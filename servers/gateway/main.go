package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/New-Era/servers/gateway/handlers"
	"github.com/New-Era/servers/gateway/models/devices"
	mgo "gopkg.in/mgo.v2"
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

	// messagingAddr := reqEnv("MESSAGESADDR")
	// summaryAddr := reqEnv("SUMMARYADDR")

	mux := http.NewServeMux()
	mux.HandleFunc("/time", handlers.TimeHandler)
	mux.HandleFunc("/device", handlerCtx.DevicesHandler)

	fmt.Printf("server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, mux))
}
