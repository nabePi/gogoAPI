// pubnubm
package contler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pubnub/go/messaging"
	"os"
	"strings"
	"time"
	//"log"
	"net/http"
)

// pub instance of messaging package.
var pub *messaging.Pubnub

// connectChannels: the conected pubnub channels, multiple channels are stored separated by comma.
var connectChannels = "kanal"

// ssl: true if the ssl is enabled else false.
var ssl bool = false

// cipher: stores the cipher key set by the user.
var cipher = ""

// uuid stores the custom uuid set by the user.
var uuid = ""

//
var publishKey = "pub-c-7611f067-063e-4ac9-be68-0fae6470e1d9"

//
var subscribeKey = "sub-c-8fbbe094-98e7-11e5-b0f3-02ee2ddab7fe"

//
var secretKey = ""

var PubnubChannel = "golife.wgs.channel"


// a boolean to capture user preference of displaying errors.
var displayError = true

type PubParam struct {
	Message string
}

func PubInit() {
	messaging.SetMaxIdleConnsPerHost(20)
	messaging.SetResumeOnReconnect(true)
	messaging.SetSubscribeTimeout(uint16(30))
	messaging.LoggingEnabled(true)

	logfileName := "pubnubMessaging.log"
	f, err := os.OpenFile(logfileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("error opening file: ", err.Error())
		fmt.Println("Logging disabled")
	} else {
		fmt.Println("Logging enabled writing to ", logfileName)
		messaging.SetLogOutput(f)
	}

	messaging.SetOrigin("pubsub.pubnub.com")

	var pubInstance = messaging.NewPubnub(publishKey, subscribeKey, secretKey, cipher, ssl, uuid)
	pub = pubInstance

	//SetupProxy()

	presenceHeartbeat := uint16(0) //askNumber16("Presence Heartbeat", true)
	pub.SetPresenceHeartbeat(presenceHeartbeat)
	//fmt.Println(fmt.Sprintf("Presence Heartbeat set to :%d", pub.GetPresenceHeartbeat()))

	presenceHeartbeatInterval := uint16(0) //askNumber16("Presence Heartbeat Interval", true)
	pub.SetPresenceHeartbeat(presenceHeartbeatInterval)
	//fmt.Println(fmt.Sprintf("Presence Heartbeat set to :%d", pub.GetPresenceHeartbeat()))

	fmt.Println("Pubnub instance initialized")

}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		//s := buf.String() // Does a complete copy of the bytes in the buffer.
		//fmt.Fprintf(w, "n = %s", s)

		//var byt []byte
		//r.Body.Read(byt)
		//		res := Response1{}
		//var dat map[string]interface{}

		res := PubParam{}

		//if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
		if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
			//			panic(err)
			//s := err.Error()
			//lr.Status = s

		} else {
			fmt.Println(res)
			go PublishRoutine(connectChannels, res.Message)
		}
		//		fmt.Println(res.Page)
	} else {

		//lr.Status = "Error parameter in body required"
	}

}

// PublishRoutine asks the user the message to send to the pubnub channel(s) and
// calls the Publish routine of the messaging package as a parallel
// process. If we have multiple pubnub channels then this method will spilt the
// channel by comma and send the message on all the pubnub channels.
func PublishRoutine(channels string, message interface{}) {
	var errorChannel = make(chan []byte)
	channelArray := strings.Split(channels, ",")

	for i := 0; i < len(channelArray); i++ {
		ch := strings.TrimSpace(channelArray[i])
		fmt.Println("Publish to channel: ", ch)
		channel := make(chan []byte)
		go pub.Publish(ch, message, channel, errorChannel)
		go handleResult(channel, errorChannel, messaging.GetNonSubscribeTimeout(), "Publish")
	}
}

func handleResult(successChannel, errorChannel chan []byte, timeoutVal uint16, action string) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(timeoutVal) * time.Second)
		timeout <- true
	}()
	for {
		select {
		case success, ok := <-successChannel:
			if !ok {
				break
			}
			if string(success) != "[]" {
				fmt.Println(fmt.Sprintf("%s Response: %s ", action, success))
				fmt.Println("")
			}
			return
		case failure, ok := <-errorChannel:
			if !ok {
				break
			}
			if string(failure) != "[]" {
				if displayError {
					fmt.Println(fmt.Sprintf("%s Error Callback: %s", action, failure))
					fmt.Println("")
				}
			}
			return
		case <-timeout:
			fmt.Println(fmt.Sprintf("%s Handler timeout after %d secs", action, timeoutVal))
			fmt.Println("")
			return
		}
	}
}
