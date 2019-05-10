# new-ERA

Create a system that acts as an intermediary between ShakeAlert earthquake warnings and the public, such as schools in Washington State. The system broadcasts earthquake warnings and aftermath effects related to earthquakes in a fast and reliable manner such that students and faculty can take appropriate actions to protect themselves. The system can be easily adopted to target different groups such as first responders and families. In addition, the system determines some baseline latency metrics between different components to offer a comparison between different solutions and evaluate the efficiency in broadcasting ShakeAlert information.
 

## Team Members

Kelley Chen | kelley97@uw.edu

Blake Eric Franzen | bfranzen@uw.edu

Jacob C. Matray | jmatray@uw.edu

Eric Wei | ericjwei@uw.edu


## API Documentation

| Endpoint |  Requests Allowed                | Description                                                                                   | Results                                                                                                                                                                                                                                                 |
|----------|----------------------|-----------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| /time       |  GET          | I want to get the current time on the server | Upon receiving a  **GET request to /time** , the handler will send a text response with the current time                                                    |
| /device       | POST | I want to create an account for my current device | Upon receiving a  **POST request to /device**,  the handler will validate the data passed to it and create a new device on success. In the event there is invalid information sent, the API will respond with an error. After validation the device will be created and information will be stored securely in the cloud. |
| /ws      |  | I want to submit an answer to a question within a game                                        | After receiving a  **POST request to /v1/trivia/{triviaID}?type=answer**,  containing a json object with the users answer,   the trivia microservice will add that to the list of answers for a specific question in the game, for later evaluation   |
| /setup       | | I want to access my statistics                                          | Upon calling the  **GET request to /v1/trivia/user/{userID}**, the user will get a json object containing the games played, points received, and number of wins |                                                 |
| /device-info/       |  | I want to send chat messages to other players                                     | Upon receiving a  **POST request to /v1/channels/{channelID}**,  the messaging microservice will insert the message body into the general chat bar in the trivia microservice                                                                                  |
| /connect       |  | I want to view chat messages sent by other players in my game                                 | When a message is posted to the message microservice it will post the message to RabbitMQ and be displayed on the trivia microservice                                                                                                                 |
| /disconnect      | | I want to view chat messages sent by other players in my game                                 | When a message is posted to the message microservice it will post the message to RabbitMQ and be displayed on the trivia microservice                                                                                                                 |
| /test      |  | I want to view chat messages sent by other players in my game                                 | When a message is posted to the message microservice it will post the message to RabbitMQ and be displayed on the trivia microservice                                                                                                   |



## Sponsors

[Sirrus7](https://www.sirrus7.com/)

## Technology and Partners

![](https://landsat.gsfc.nasa.gov/wp-content/uploads/2013/09/USGS_logo_green.png)
![](https://pbs.twimg.com/profile_images/692813728446722048/7kg5YJ6F_400x400.png)
![](https://freeicons.io/laravel/public/uploads/icons/png/18181230061536126577-128.png)
![](https://rallen.berkeley.edu/img/ShakeAlertlogo.png)
