package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messageId uint64

type PayloadGenerator func(i int) string

func GenerateMessageBaseValue() {
	rand.Seed(time.Now().UnixMicro())
	messageId = randomSource.Uint64()
}

func defaultPayloadGen() PayloadGenerator {
	return func(i int) string {
		return fmt.Sprintf("this is msg #%d!", i)
	}
}

func constantPayloadGenerator(payload string) PayloadGenerator {
	return func(i int) string {
		return payload
	}
}

func filePayloadGenerator(filepath string) PayloadGenerator {
	inputPath := strings.Replace(filepath, "@", "", 1)
	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("error reading payload file: %v\n", err)
		os.Exit(1)
	}
	return func(i int) string {
		return string(content)
	}
}

type Worker struct {
	WorkerId             int
	BrokerUrl            string
	Username             string
	Password             string
	SkipTLSVerification  bool
	NumberOfMessages     int
	PayloadGenerator     PayloadGenerator
	Timeout              time.Duration
	Retained             bool
	PublisherQoS         byte
	SubscriberQoS        byte
	CA                   []byte
	Cert                 []byte
	Key                  []byte
	PauseBetweenMessages time.Duration
}

func setSkipTLS(o *mqtt.ClientOptions) {
	oldTLSCfg := o.TLSConfig
	if oldTLSCfg == nil {
		oldTLSCfg = &tls.Config{}
	}
	oldTLSCfg.InsecureSkipVerify = true
	o.SetTLSConfig(oldTLSCfg)
}
func NewTlsConfig2() *tls.Config {
	certpool := x509.NewCertPool()
	ca := "-----BEGIN CERTIFICATE-----\nMIIDXDCCAkSgAwIBAgITBcMGZIDrc0uW9GCK0Q8YrfBRgjANBgkqhkiG9w0BAQsF\nADA2MRkwFwYDVQQKExBrb3Jld2lyZWxlc3MuY29tMRkwFwYDVQQDExBrb3Jld2ly\nZWxlc3MuY29tMB4XDTIyMDkyMjEzMjAwOFoXDTMyMDkxOTEzMjAwN1owNjEZMBcG\nA1UEChMQa29yZXdpcmVsZXNzLmNvbTEZMBcGA1UEAxMQa29yZXdpcmVsZXNzLmNv\nbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAIzyuJG1lqM38Q9uGaEI\nm+15LHqcl9HWvxxOpuJxVB4+p9mcQiGqXf4K8Oyz5VGrFMNtZ1qC4jwrBJq6qZwx\nv/i37JFEHOpwbFEEz8zeSxIG3OCGsD1fnafwjD+kVP+ObaI6oN7cqEtRgbPDu/Qz\nW9eH//6o11H2q/CVXgiDlLTOqQjPNefoqmwy//xDzUV5nfBs21h/c2VFHwkPgJbo\n6CvFgdmi25JS4nQnB8qNlhilzuRfjwxUSjeMuuaGSxLmojCgTFrkCOokjAGs7KzU\n6vKT7juRMk0QA8V6/7Clum+p5pEKXdBxtdQ3KJrWUqENOi62VCQvu6gtTSfpSnV0\n5hUCAwEAAaNjMGEwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYD\nVR0OBBYEFNnVdulRaYA5kosaQuZOHQH98CEZMB8GA1UdIwQYMBaAFNnVdulRaYA5\nkosaQuZOHQH98CEZMA0GCSqGSIb3DQEBCwUAA4IBAQB09XZmi3HTifsfF2OMOupo\nlKxys5LRULGDBVjvwEd+6PjWP0yXRC2bcfITPN3dkZF/VJsGtZZVFslOpkErQQhs\nkh8+uZ9UcmqKhwi49UcFOksRMl37o7f0Wp9o71AsOIybBoMe8BDWLojgYMUjgJym\nIIrHa0Xer+n33KYLBmPplfWSmrD7KBS+Gejg7fuVa3EjfniX4/oVeW1OOxsNWBF5\n0GMxt8P73j8VWXuO1f068MdwrZsIVxPo6+NpnU4R5JtsgYsR+zNj/GKJLg0M2EQ8\nJ4GJmao7FCScH1j6nMP1BliDJEKtfxD2QhwPX2FFZHOgRwBombmF/BymBtyWRjAn\n-----END CERTIFICATE-----\n"
	bytesCa := []byte(ca)
	certpool.AppendCertsFromPEM(bytesCa)
	return &tls.Config{
		RootCAs: certpool,
	}
}
func NewTLSConfig(ca, certificate, privkey []byte) (*tls.Config, error) {
	// Import trusted certificates from CA
	certpool := x509.NewCertPool()
	ok := certpool.AppendCertsFromPEM(ca)

	if !ok {
		return nil, fmt.Errorf("CA is invalid")
	}

	// Import client certificate/key pair
	_, err := tls.X509KeyPair(certificate, privkey)
	if err != nil {
		return nil, err
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		//Certificates: []tls.Certificate{cert},
		Certificates: nil,
	}, nil
}

func (w *Worker) Run(ctx context.Context) {
	verboseLogger.Printf("[%d] initializing\n", w.WorkerId)

	queue := make(chan [2]string)
	cid := w.WorkerId
	_ = randomSource.Int31()

	_, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	topicName := fmt.Sprintf(topicNameTemplate, w.WorkerId)
	subscriberClientId := fmt.Sprintf(subscriberClientIdTemplate, w.WorkerId)
	publisherClientId := fmt.Sprintf(publisherClientIdTemplate, w.WorkerId)
	verboseLogger.Printf("[%d] topic=%s subscriberClientId=%s publisherClientId=%s\n", cid, topicName, subscriberClientId, publisherClientId)
	password := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJLb3JlV2lyZWxlc3MiLCJleHAiOjM1NjEzNDc3NzMsImlhdCI6MTY2Nzg5MTc3M30.o003E_EpJf5uND3Gmpbw5EcGOUz4eoW28pecsnRbiNSsRqC138yjb1j3GfoIB1tKNFtj5pqGYyh0Unx4IV_ubvwGM1RbPFM3LK4vjnReJpH36KOFDwqHnofxEk3RZrcmjX6q3hprPgw_0YEEG2nAew21x7Dy5AwJXbLE_xBfjHGVoJUijUub0WZici1QGG05bpu8CcAMgDzBqmc-KHq1On1gyny0MjFgT3l30UVIxj_x_g0nQmRQVcEDaIoRuMNm5HDZ0DL9Aw_Zl6KgcH59CfUpDkwylkseYt5W7R3n9lWQ7AqLvawKb6FbNwrdZTnILSYz005W_eh7M1VL4ml6Eg"
	tlsConfig := NewTlsConfig2()
	publisherOptions := mqtt.NewClientOptions().SetClientID(publisherClientId).SetUsername("unused").SetPassword(password).AddBroker(w.BrokerUrl)

	subscriberOptions := mqtt.NewClientOptions().SetClientID(subscriberClientId).SetUsername("unused").SetPassword(password).AddBroker(w.BrokerUrl)
	publisherOptions.SetTLSConfig(tlsConfig)
	subscriberOptions.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		queue <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	if len(w.CA) > 0 || len(w.Key) > 0 {
		tlsConfig, err := NewTLSConfig(w.CA, w.Cert, w.Key)
		if err != nil {
			panic(err)
		}
		publisherOptions.SetTLSConfig(tlsConfig)
	}

	if w.SkipTLSVerification {
		setSkipTLS(publisherOptions)
		setSkipTLS(subscriberOptions)
	}

	// subscriber := mqtt.NewClient(subscriberOptions)

	// verboseLogger.Printf("[%d] connecting subscriber\n", w.WorkerId)
	// if token := subscriber.Connect(); token.WaitTimeout(w.Timeout) && token.Error() != nil {
	// 	resultChan <- Result{
	// 		WorkerId:     w.WorkerId,
	// 		Event:        ConnectFailedEvent,
	// 		Error:        true,
	// 		ErrorMessage: token.Error(),
	// 	}

	// 	return
	// }

	// defer func() {
	// 	verboseLogger.Printf("[%d] unsubscribe\n", w.WorkerId)

	// 	if token := subscriber.Unsubscribe(topicName); token.WaitTimeout(w.Timeout) && token.Error() != nil {
	// 		fmt.Printf("failed to unsubscribe: %v\n", token.Error())
	// 	}

	// 	subscriber.Disconnect(5)
	// }()

	// verboseLogger.Printf("[%d] subscribing to topic\n", w.WorkerId)
	// if token := subscriber.Subscribe(topicName, w.SubscriberQoS, nil); token.WaitTimeout(w.Timeout) && token.Error() != nil {
	// 	resultChan <- Result{
	// 		WorkerId:     w.WorkerId,
	// 		Event:        SubscribeFailedEvent,
	// 		Error:        true,
	// 		ErrorMessage: token.Error(),
	// 	}

	// 	return
	// }

	publisher := mqtt.NewClient(publisherOptions)
	verboseLogger.Printf("[%d] connecting publisher\n", w.WorkerId)
	if token := publisher.Connect(); token.WaitTimeout(w.Timeout) && token.Error() != nil {
		resultChan <- Result{
			WorkerId:     w.WorkerId,
			Event:        ConnectFailedEvent,
			Error:        true,
			ErrorMessage: token.Error(),
		}
		return
	}

	verboseLogger.Printf("[%d] starting control loop %s\n", w.WorkerId, topicName)

	receivedCount := 0
	publishedCount := 0

	t0 := time.Now()
	for i := 0; i < w.NumberOfMessages; i++ {
		text := fmt.Sprintf("{\"id\":%d,\"time\":%d}", messageId, time.Now().UTC().UnixMilli())
		atomic.AddUint64(&messageId, 1)
		token := publisher.Publish(topicName, w.PublisherQoS, w.Retained, text)
		published := token.WaitTimeout(w.Timeout)
		if published {
			publishedCount++
		}
		time.Sleep(w.PauseBetweenMessages)
	}
	publisher.Disconnect(5)

	publishTime := time.Since(t0)
	verboseLogger.Printf("[%d] all messages published\n", w.WorkerId)
	resultChan <- Result{
		WorkerId:          w.WorkerId,
		Event:             CompletedEvent,
		PublishTime:       publishTime,
		ReceiveTime:       time.Since(t0),
		MessagesReceived:  receivedCount,
		MessagesPublished: publishedCount,
	}

	verboseLogger.Printf("[%d] worker finished\n", w.WorkerId)
}
