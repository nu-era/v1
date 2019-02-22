package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

)

const headerUser = "X-User"

var sessStore *sessions.RedisStore
var sessionKey string

// main entry point for the server
func main() {
	sessionKey = os.Getenv("SESSIONKEY")
	if len(sessionKey) == 0 {
		log.Fatal("please set session key")
	}
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	tlscert := os.Getenv("TLSCERT")
	tlskey := os.Getenv("TLSKEY")
	if len(tlskey) == 0 && len(tlscert) == 0 {
		log.Fatal("please set TLSKEY and TLSCERT")
	}

	// // construct rediserver
	// redisAddr := reqEnv("REDISADDR")
	// rClient := redis.NewClient(&redis.Options{
	// 	Addr:     redisAddr,
	// 	Password: "",
	// 	DB:       0,
	// })

	// // construct mysql server
	// dsn := fmt.Sprintf("root:%s@tcp(mysql:3306)/mysql",
	// 	os.Getenv("MYSQL_ROOT_PASSWORD"))
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	fmt.Printf("error opening database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer db.Close()

	// // connect to RabbitMq
	// conn, err := amqp.Dial("amqp://guest:guest@rabbitMq:5672/")
	// failOnError(err, "Failed to connect to RabbitMq")
	// defer conn.Close()

	// ch, err := conn.Channel()
	// failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"messages", // name
	// 	true,       // durable
	// 	false,      // delete when unused
	// 	false,      // exclusive
	// 	false,      // no-wait
	// 	nil,        // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// err = ch.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	// failOnError(err, "Failed to set QoS")

	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	false,  // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// failOnError(err, "Failed to register a consumer")

	// forever := make(chan []byte)

	// go func() {
	// 	for m := range msgs {
	// 		forever <- m.Body
	// 		m.Ack(false)
	// 	}
	// }()

	// sessStore = sessions.NewRedisStore(rClient, time.Duration(600)*time.Second)
	// userStore := users.NewMySQLStore(db)
	// searchIndex := indexes.NewTrie()
	// userStore.Populate(searchIndex)
	// socketStore := handlers.NewSocketStore(forever)
	hctx := handlers.NewHandlerContext(sessionKey, sessStore, userStore, searchIndex, socketStore)

	// messagingAddr := reqEnv("MESSAGESADDR")
	// summaryAddr := reqEnv("SUMMARYADDR")

	mux := http.NewServeMux()
	// mux.HandleFunc("/v1/users", hctx.UsersHandler)
	// mux.HandleFunc("/v1/users/", hctx.SpecificUserHandler)
	// mux.HandleFunc("/v1/sessions", hctx.SessionsHandler)
	// mux.HandleFunc("/v1/sessions/", hctx.SpecificSessionHandler)
	// mux.HandleFunc("/ws", hctx.WebSocketConnectionHandler)
	// mux.Handle("/v1/summary", NewServiceProxy(summaryAddr, rClient))
	// mux.Handle("/v1/channels", NewServiceProxy(messagingAddr, rClient))
	// mux.Handle("/v1/channels/", NewServiceProxy(messagingAddr, rClient))
	// mux.Handle("/v1/messages/", NewServiceProxy(messagingAddr, rClient))
	wrappedMux := handlers.NewCorsHandler(mux)

	fmt.Printf("server is listening at https://%s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, wrappedMux))
}

// //NewServiceProxy returns a new ReverseProxy for a micorservice given a comma-delimited list
// //of network addresses
// func NewServiceProxy(addrs string, rc *redis.Client) *httputil.ReverseProxy {
// 	// _addrs, _ := rc.Get(addrs).Result()
// 	// splitAddrs := strings.Split(_addrs, ",")
// 	splitAddrs := strings.Split(addrs, ",")
// 	nextAddr := 0
// 	mx := sync.Mutex{}

// 	return &httputil.ReverseProxy{
// 		Director: func(r *http.Request) {
// 			r.Header.Del("X-User")
// 			user, err := getUser(r)
// 			if err == nil {
// 				u, err := json.Marshal(user)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				r.Header.Set("X-User", string(u))
// 			}
// 			mx.Lock()
// 			r.URL.Scheme = "http"
// 			r.Host = splitAddrs[nextAddr]
// 			r.URL.Host = splitAddrs[nextAddr]
// 			nextAddr = (nextAddr + 1) % len(splitAddrs)
// 			mx.Unlock()
// 		},
// 	}
// }

// //reqEnv helper for getting list of network addresses
// func reqEnv(name string) string {
// 	val := os.Getenv(name)
// 	if len(val) == 0 {
// 		log.Fatalf("please set the %s environment variable", name)
// 	}
// 	return val
// }

// //getUser gets the user of the current sessionState
// func getUser(r *http.Request) (*users.User, error) {
// 	sessionState := &handlers.SessionState{}
// 	_, err := sessions.GetState(r, sessionKey, sessStore, sessionState)
// 	if err != nil {
// 		fmt.Printf("problem with session %v", err)
// 	}
// 	return sessionState.User, nil
// }

// //failOnError error messages
// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }