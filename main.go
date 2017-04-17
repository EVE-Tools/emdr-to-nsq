package main

import (
	"bytes"
	"compress/zlib"

	"regexp"
	"runtime"
	"time"

	"github.com/EVE-Tools/emdr-to-nsq/lib/emds"
	"github.com/EVE-Tools/emdr-to-nsq/lib/messageProcessing"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/kelseyhightower/envconfig"
	"github.com/mailru/easyjson"
	"github.com/nsqio/go-nsq"
	"github.com/pebbe/zmq4"
)

// Config holds the application's configuration info from the environment.
type Config struct {
	LogLevel         string         `default:"info" envconfig:"log_level"`
	EMDRRelayURL     string         `default:"tcp://relay-eu-germany-1.eve-emdr.com:8050" envconfig:"emdr_relay_url"`
	NSQURL           string         `default:"nsqd:4150" envconfig:"nsq_url"`
	GeneratorName    string         `envconfig:"generator_name"`
	GeneratorVersion string         `envconfig:"generator_version"`
	CachePath        string         `default:"cache.db" envconfig:"cache_path"`
	NameRegex        *regexp.Regexp `ignored:"true"`
	VersionRegex     *regexp.Regexp `ignored:"true"`
}

// Stores main configuration
var config Config

// Channel for messages to be sent to NSQ
var nsqUpstream = make(chan []byte, 10000)

func main() {
	// Load config and connect to queues
	loadConfig()
	connectToBolt()
	connectToEMDR()
	connectToNSQ()

	// Terminate this goroutine, crash if all other goroutines exited
	runtime.Goexit()
}

// Load configuration from environment and compile regexps
func loadConfig() {
	envconfig.MustProcess("EMDR_TO_NSQ", &config)

	logLevel, err := log.ParseLevel(config.LogLevel)

	if err != nil {
		panic(err)
	}

	log.SetLevel(logLevel)

	if config.GeneratorName != "" {
		config.NameRegex = regexp.MustCompile(config.GeneratorName)
	}

	if config.GeneratorVersion != "" {
		config.VersionRegex = regexp.MustCompile(config.GeneratorVersion)
	}

	log.Info("\n    ________  _______  ____     __           _   _______ ____ \n   / ____/  |/  / __ \\/ __ \\   / /_____     / | / / ___// __ \\\n  / __/ / /|_/ / / / / /_/ /  / __/ __ \\   /  |/ /\\__ \\/ / / /\n / /___/ /  / / /_/ / _, _/  / /_/ /_/ /  / /|  /___/ / /_/ /\n/_____/_/  /_/_____/_/ |_|   \\__/\\____/  /_/ |_//____/\\___\\_\\ v0.1")
	log.Debugf("Config: %q", config)
}

// Connect to local boltdb
func connectToBolt() {
	db, err := bolt.Open(config.CachePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	messageProcessing.Initialize(db)
}

// Connect to EMDR and start listening for messages
func connectToEMDR() {
	// create conn and pass to go listenToEMDR
	connection, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}

	err = connection.SetSubscribe("")
	if err != nil {
		panic(err)
	}

	err = connection.Connect(config.EMDRRelayURL)
	if err != nil {
		panic(err)
	}

	//  Ensure subscriber connection has time to complete
	time.Sleep(time.Second)

	go listenToEMDR(connection)
}

// Connect to NSQ and start processing messages
func connectToNSQ() {
	// create conn and pass to go pushToNSQ
	nsqConfig := nsq.NewConfig()
	nsqConfig.Snappy = false
	producer, err := nsq.NewProducer(config.NSQURL, nsqConfig)
	if err != nil {
		panic(err)
	}

	go pushToNSQ(producer)
}

// Push messages received to NSQ
func pushToNSQ(producer *nsq.Producer) {
	defer producer.Stop()

	for {
		message := <-nsqUpstream
		err := producer.Publish("orders", message)
		if err != nil {
			log.Warnf("Could not publish message: ", err.Error())
		}
	}
}

// Listen for EMDR messages
func listenToEMDR(connection *zmq4.Socket) {

	defer connection.Close()

	for {
		message, err := connection.RecvBytes(0)
		if err != nil {
			log.Errorf("Recv error: %s", err.Error())
			continue
		}

		go processMessage(message)

	}
}

func processMessage(rawMessage []byte) {
	message, err := decompress(rawMessage)
	if err != nil {
		log.Warnf("Failed to decompress message: %s", err.Error())
		return
	}

	rowsets, err := messageProcessing.FilterMessage(config.NameRegex, config.VersionRegex, message)
	if err != nil {
		return
	}
	if len(rowsets) > 0 {
		pushToOutbox(rowsets)
	}
}

// Marshals and sends rowsets to NSQ goroutine.
func pushToOutbox(rowsets []emds.Rowset) error {
	for _, rowset := range rowsets {
		json, err := easyjson.Marshal(rowset)
		if err != nil {
			return err
		}

		nsqUpstream <- json
	}

	return nil
}

// Decompresses message from EMDR.
func decompress(rawMessage []byte) ([]byte, error) {
	// Decompress
	rawMessageBuffer := bytes.NewReader(rawMessage)
	rawMessageReader, err := zlib.NewReader(rawMessageBuffer)
	if err != nil {
		return nil, err
	}

	messageBuffer := new(bytes.Buffer)

	_, err = messageBuffer.ReadFrom(rawMessageReader)
	if err != nil {
		return nil, err
	}
	rawMessageReader.Close()

	return messageBuffer.Bytes(), nil
}
