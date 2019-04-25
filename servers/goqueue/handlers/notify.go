package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/akavel/polyclip-go"
	"github.com/streadway/amqp"
	"strconv"
	"strings"
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
			- make contour based on polygon data from API
			- filter devices based on location relative to MMI 4 polygon if present
			- publish message to queue
	*/
}

// makeContour takes in lat/longitude point data from the ShakeAlert API
// and generates/returns a polyclip contour object
func makeContour(data string) *polyclip.Contour {
	points := strings.Split(data, " ")
	contour := &polyclip.Contour{}
	for _, p := range points {
		coords := strings.Split(p, ",")
		x, _ := strconv.ParseFloat(coords[0], 64)
		y, _ := strconv.ParseFloat(coords[1], 64)
		contour.Add(polyclip.Point{
			X: x,
			Y: y,
		})
	}
	return contour
}