package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var executions int = 0
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func sendMessage(w http.ResponseWriter, r *http.Request, snsarn string, arn string, externalID string, region string, maxMessagesPerMinute int) {

	keys, ok := r.URL.Query()["key"]
	var key string
	if !ok || len(keys[0]) < 1 {
		key = ""
	} else {
		key = keys[0]
	}

	sess := session.Must(session.NewSession())
	conf := createConfig(arn, externalID, region, sess)

	if executions < maxMessagesPerMinute {
		svc := sns.New(sess, &conf)
		params := &sns.PublishInput{
			Message:  aws.String(key),
			TopicArn: aws.String(snsarn),
		}
		resp, error := svc.Publish(params)
		if error != nil {
			Error.Println("Publish failed", error)
		} else {
			Info.Println("Message sent. Number of executions: %d. SNS ARN: %s, Region: %s, response: %s", executions+1, snsarn, arn, region, resp)
		}
		executions += 1
	} else {
		Warning.Println("Message not sent, because the per minute limit has been reached. Number of executions > %d", executions-1)
	}
}

func createConfig(arn string, externalID string, region string, sess *session.Session) aws.Config {

	conf := aws.Config{Region: aws.String(region)}
	if arn != "" {
		// if ARN flag is passed in, we need to be able ot assume role here
		var creds *credentials.Credentials
		if externalID != "" {
			// If externalID flag is passed, we need to include it in credentials struct
			creds = stscreds.NewCredentials(sess, arn, func(p *stscreds.AssumeRoleProvider) {
				p.ExternalID = &externalID
			})
		} else {
			creds = stscreds.NewCredentials(sess, arn, func(p *stscreds.AssumeRoleProvider) {})
		}
		conf.Credentials = creds
	}
	return conf
}

func createResetTicker() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for t := range ticker.C {
			Trace.Println("Tick at", t, executions)
			executions = 0
		}
	}()
}

func main() {

	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	var arn string
	var externalID string
	var region string
	var snsarn string
	var maxMessagesPerMinute int
	var port int
	const (
		defaultARN    = ""
		arnUsage      = "The ARN of the role you need to assume"
		defaultExtID  = ""
		extIDUsage    = "The ExternalID constraint, if applicable for the role you need to assume"
		defaultRegion = ""
		regionUsage   = "The region of the SNS topic (mandatory)"
		defaultSNSARN = ""
		SNSARNUsage   = "The ARN of the receiver SNS topic (mandatory)"
		MMPMUsage     = "The maximum number of messages allowed per minute"
		defaultMMPM   = 20
		portUsage     = "The listening port for the application"
		defaultPort   = 8080
	)
	flag.StringVar(&arn, "arn", defaultARN, arnUsage)
	flag.StringVar(&externalID, "extid", defaultExtID, extIDUsage)
	flag.StringVar(&region, "region", defaultRegion, regionUsage)
	flag.StringVar(&snsarn, "snsarn", defaultSNSARN, SNSARNUsage)
	flag.IntVar(&maxMessagesPerMinute, "maxMessagesPerMinute", defaultMMPM, MMPMUsage)
	flag.IntVar(&port, "port", defaultPort, portUsage)
	flag.Parse()
	if snsarn == "" || region == "" {
		Error.Println("Please supply the mandatory parameters, snsarn and region", snsarn, region)
		os.Exit(1)
	}

	createResetTicker()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sendMessage(w, r, snsarn, arn, externalID, region, maxMessagesPerMinute)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
