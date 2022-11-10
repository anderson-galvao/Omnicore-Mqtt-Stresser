package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/RacoWireless/Omnicore-Mqtt-Stresser/model"
)

var (
	SummaryChannelData = make(chan Summary, 1)
	resultChan         = make(chan Result)
	stopWaitLoop       = false
	randomSource       = rand.New(rand.NewSource(time.Now().UnixNano()))

	errorLogger   = log.New(os.Stderr, "ERROR: ", log.Lmicroseconds|log.Ltime|log.Lshortfile)
	verboseLogger = log.New(os.Stderr, "DEBUG: ", log.Lmicroseconds|log.Ltime|log.Lshortfile)
)

type Result struct {
	WorkerId          int
	Event             string
	PublishTime       time.Duration
	ReceiveTime       time.Duration
	MessagesReceived  int
	MessagesPublished int
	Error             bool
	ErrorMessage      error
}

type TimeoutError interface {
	Timeout() bool
	Error() string
}

func parseQosLevels(qos int, role string) (byte, error) {
	if qos < 0 || qos > 2 {
		return 0, fmt.Errorf("%q is an invalid QoS level for %s. Valid levels are 0, 1 and 2", qos, role)
	}
	return byte(qos), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// An error is returned of the given TLS configuration is invalid.
func validateTLSFiles(argCafile, argKey, argCert string) error {
	if len(argCafile) > 0 {
		if !fileExists(argCafile) {
			return fmt.Errorf("CA file %q does not exist", argCafile)
		}
	}
	if len(argKey) > 0 {
		if !fileExists(argKey) {
			return fmt.Errorf("key file %q does not exist", argKey)
		}
	}
	if len(argCert) > 0 {
		if !fileExists(argCert) {
			return fmt.Errorf("cert file %q does not exist", argCert)
		}
	}

	if len(argKey) > 0 && len(argCert) < 1 {
		return fmt.Errorf("A key file is specified but no certificate file")
	}

	if len(argKey) < 1 && len(argCert) > 0 {
		return fmt.Errorf("A cert file is specified but no key file")
	}
	return nil
}

// loadTLSFile loads the given file. If the filename is empty neither data nor an error is returned.
func loadTLSFile(fileName string) ([]byte, error) {
	if len(fileName) > 0 {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS file: %q: %w", fileName, err)
		}
		return data, nil
	}
	return nil, nil
}

func (d *StresserService) ExecuteStresser(Arguments model.Stresser, tenant string) error {
	// flag.Parse()

	// if flag.NFlag() < 1 || *argHelp {
	// 	flag.Usage()
	// 	if *argHelp {
	// 		return errors.New("Arg Help Not Found")
	// 	}
	// 	return errors.New("Arg Help Not Found")
	// }
	tenantTemplate := fmt.Sprintf("tenants/%s", tenant)
	subscriberClientIdTemplate := tenantTemplate + "/locations/us-central1/registries/KoreWireless/devices/Stresser%d"
	publisherClientIdTemplate := tenantTemplate + "/locations/us-central1/registries/KoreWireless/devices/Stresser%d"
	topicNameTemplate := tenantTemplate + "/registries/KoreWireless/devices/Stresser%d/events"
	argNumClients := Arguments.Clients
	argNumMessages := Arguments.Messages
	argConstantPayload := ""
	argTimeout := "5s"
	argGlobalTimeout := "20s"
	argRampUpSize := 100
	argRampUpDelay := "500ms"
	argBrokerUrl := d.BrokerUrl
	argUsername := "unused"
	argPassword := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJteS1pb3QtMzU2MzA1IiwiZXhwIjoxNjk1NjY3OTc0LCJpYXQiOjE2NjQxMTAzNzR9.T_kzjb2mQVtF_0J9zY7QuJiY8z5sd8-VNN8XW06xo1CGQvpjYnOcfVs0tfh6t8VWDZq5PndcbNTNCybZbJd4Dhzxw_Rz-6PJoFqe9HisIl7xyRNanxzVEeeBE-3SSmJRSPTGYjx6VHZU2xRYCNmXSi0UdLPi6P43-TdK3gPZDR57CJQbbGUdVSotVAz9tbETNBdthZK6tpw8o8EgKpsBfKKOzNmXYAtt9wHuoPSI_HlFSviMMEEYZuC8Ss3xJ6nGWJuQEY6G4epsrnjxneT3fHGcjflI-if4FmdRmxmcvCQBrZd2UGvylJTK96Ir3WQfcJbQdT2n9Fc7VVifYR3Lzw"
	argLogLevel := 0
	argProfileCpu := ""
	argProfileMem := ""
	argHideProgress := false
	argRetain := false
	argPublisherQoS := 1
	argSubscriberQoS := 1
	argSkipTLSVerification := false
	argCafile := ""
	argKey := ""
	argCert := ""
	argPauseBetweenMessages := "0s"
	argTopicBasePath := ""
	if argProfileCpu != "" {
		f, err := os.Create(argProfileCpu)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create CPU profile: %s\n", err)
			return err
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start CPU profile: %s\n", err)
			return err
		}
	}

	num := argNumMessages
	username := argUsername
	password := argPassword

	actionTimeout, err := time.ParseDuration(argTimeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse '--timeout': %q is not a valid duration string. See https://golang.org/pkg/time/#ParseDuration for valid duration strings\n", argTimeout)
		return err
	}

	verboseLogger.SetOutput(ioutil.Discard)
	errorLogger.SetOutput(ioutil.Discard)

	if argLogLevel == 1 || argLogLevel == 3 {
		errorLogger.SetOutput(os.Stderr)
	}

	if argLogLevel == 2 || argLogLevel == 3 {
		verboseLogger.SetOutput(os.Stderr)
	}

	if argBrokerUrl == "" {
		fmt.Fprintln(os.Stderr, "'--broker' is empty. Abort.")
		return err
	}

	if len(argTopicBasePath) > 0 {
		topicNameTemplate = strings.Replace(topicNameTemplate, "internal/mqtt-stresser", argTopicBasePath, 1)
	}

	payloadGenerator := defaultPayloadGen()
	if len(argConstantPayload) > 0 {
		if strings.HasPrefix(argConstantPayload, "@") {
			verboseLogger.Printf("Set constant payload from file %s\n", argConstantPayload)
			payloadGenerator = filePayloadGenerator(argConstantPayload)
		} else {
			verboseLogger.Printf("Set constant payload to %s\n", argConstantPayload)
			payloadGenerator = constantPayloadGenerator(argConstantPayload)
		}
	}

	var publisherQoS, subscriberQoS byte

	if lvl, err := parseQosLevels(argPublisherQoS, "publisher"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	} else {
		publisherQoS = lvl
	}

	if lvl, err := parseQosLevels(argSubscriberQoS, "subscriber"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	} else {
		subscriberQoS = lvl
	}

	var ca, cert, key []byte
	if err := validateTLSFiles(argCafile, argKey, argCert); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	ca, err = loadTLSFile(argCafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	cert, err = loadTLSFile(argCert)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	key, err = loadTLSFile(argKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	rampUpDelay, _ := time.ParseDuration(argRampUpDelay)
	rampUpSize := argRampUpSize

	if rampUpSize < 0 {
		rampUpSize = 100
	}

	resultChan = make(chan Result, argNumClients*argNumMessages)

	globalTimeout, err := time.ParseDuration(argGlobalTimeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parse '--global-timeout': %q is not a valid duration string. See https://golang.org/pkg/time/#ParseDuration for valid duration strings\n", argGlobalTimeout)
		return err
	}
	testCtx, cancelFunc := context.WithTimeout(context.Background(), globalTimeout)
	defer cancelFunc()
	pauseBetweenMessages, err := time.ParseDuration(argPauseBetweenMessages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parse '--pause-between-messages': %q is not a valid duration string. See https://golang.org/pkg/time/#ParseDuration for valid duration strings\n", argPauseBetweenMessages)
		return err
	}
	stopStartLoop := false
	for cid := 0; cid < argNumClients && !stopStartLoop; cid++ {

		if cid%rampUpSize == 0 && cid > 0 {
			fmt.Printf("%d worker started - waiting %s\n", cid, rampUpDelay)
			select {
			case <-time.NewTimer(rampUpDelay).C:
			case s := <-signalChan:
				fmt.Printf("Got signal %s. Cancel test.\n", s.String())
				cancelFunc()
				stopStartLoop = true
			}
		}

		go (&Worker{
			tenant:                     tenant,
			publisherClientIdTemplate:  publisherClientIdTemplate,
			subscriberClientIdTemplate: subscriberClientIdTemplate,
			topicNameTemplate:          topicNameTemplate,
			WorkerId:                   cid,
			BrokerUrl:                  argBrokerUrl,
			Username:                   username,
			Password:                   password,
			SkipTLSVerification:        argSkipTLSVerification,
			NumberOfMessages:           num,
			PayloadGenerator:           payloadGenerator,
			Timeout:                    actionTimeout,
			Retained:                   argRetain,
			PublisherQoS:               publisherQoS,
			SubscriberQoS:              subscriberQoS,
			CA:                         ca,
			Cert:                       cert,
			Key:                        key,
			PauseBetweenMessages:       pauseBetweenMessages,
		}).Run(testCtx)
	}
	fmt.Printf("%d worker started\n", argNumClients)
	finEvents := 0

	results := make([]Result, argNumClients)
	for finEvents < argNumClients && !stopWaitLoop {
		select {
		case msg := <-resultChan:
			results[msg.WorkerId] = msg

			if msg.Event == CompletedEvent || msg.Error {
				finEvents++
				verboseLogger.Printf("%d/%d events received\n", finEvents, argNumClients)
			}

			if msg.Error {
				errorLogger.Println(msg)
			}

			if !argHideProgress {
				if msg.Event == ProgressReportEvent {
					fmt.Print(".")
				}

				if msg.Error {
					fmt.Print("E")
				}
			}

		case <-testCtx.Done():
			switch testCtx.Err().(type) {
			case TimeoutError:
				fmt.Println("Test timeout. Wait 5s to allow disconnection of clients.")
			default:
				fmt.Println("Test canceled. Wait 5s to allow disconnection of clients.")
			}
			time.Sleep(5 * time.Second)
			stopWaitLoop = true
		case s := <-signalChan:
			fmt.Printf("Got signal %s. Cancel test.\n", s.String())
			cancelFunc()
			stopWaitLoop = true
		}
	}

	summary, err := buildSummary(argNumClients, num, results, tenant)
	if err != nil {
		return err
	} else {
		printSummary(summary)
		SummaryChannelData <- summary
	}
	if argProfileMem != "" {
		f, err := os.Create(argProfileMem)

		if err != nil {
			fmt.Printf("Failed to create memory profile: %s\n", err)
		}

		runtime.GC() // get up-to-date statistics

		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Printf("Failed to write memory profile: %s\n", err)
		}
		f.Close()
	}
	pprof.StopCPUProfile()
	return nil
}
