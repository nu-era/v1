package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/akavel/polyclip-go"
	"github.com/streadway/amqp"
)

// map for shaking intensity based on MMI
var shaking = map[int]string{
	4:  "Light",
	5:  "Moderate",
	6:  "Strong",
	7:  "Very Strong",
	8:  "Severe",
	9:  "Violent",
	10: "Extreme",
}

// PublishData takes the input data and publishes it to rabbitmq
// for consumers to parse and send to clients
func (ctx *QueueContext) PublishData(data interface{}, name string) {
	body, _ := json.Marshal(data)

	queue, err := ctx.Channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error declaring queue, %v", err)
	}

	err = ctx.Channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		fmt.Printf("error publish to queue, %v", err)
	}
}

// getDevices is a function that filters through devices when an alert is received
// from the ShakeAlert API to determine what devices need to be notified
func (ctx *QueueContext) getDevices() {
	// what are we getting for parameters?
	// possibly polygon we use to filter lat/long of devices on

}

func (ctx *QueueContext) Routine() {
	/*TODO:
	- listen for data coming in from ShakeAlert API
		- publish data to receive queue?
			- read off of receive queue
			- filter devices based on location relative to MMI 4 polygon if present
			- publish message to queue
	-
	*/
}

// makePolygon takes in lat/longitude point data from the ShakeAlert API
// and generates/returns a polyclip polygon object
func makePolygon(data interface{}) *polyclip.Polygon {
	return nil
}
