package main
import (
  "time"
  "flag"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
  "github.com/aws/aws-sdk-go/aws/session"
  "net/http"
  "fmt"
  "os"
  "github.com/aws/aws-sdk-go/service/sns"
)

var executions int     = 0

func sendMessage(w http.ResponseWriter, r *http.Request, snsarn string, arn string, externalID string, region string) {

	sess := session.Must(session.NewSession())
	conf := createConfig(arn, externalID, region, sess)

        if  executions < 1 {
            fmt.Println("Message sent. Number of executions:", executions+1, snsarn, arn, region, externalID)
            svc := sns.New(sess, &conf)
            params := &sns.PublishInput{
                Message: aws.String("message"), 
                TopicArn: aws.String(snsarn),  
            }
            resp, error := svc.Publish(params)
            if error != nil {
            fmt.Println("Publish failed", error)
            }
            executions += 1
            w.Write([]byte(fmt.Sprintf("%v", resp)))
        } else {
            w.Write([]byte(fmt.Sprintf("Message not sent. Number of executions > ", executions)))
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

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
      return value
  }
  return fallback
}

func createResetTicker() {
  ticker := time.NewTicker(10 * time.Second)
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t, executions)
            executions = 0
        }
    }()
}

func main() {

var arn string
var externalID string
var region string
var snsarn string
const (
  
  defaultARN     = ""
  arnUsage       = "The ARN of the role you need to assume"
  defaultExtID   = ""
  extIDUsage     = "The ExternalID constraint, if applicable for the role you need to assume"
  defaultRegion  = ""
  regionUsage    = "The region of the SNS topic (mandatory)"
  defaultSNSARN  = "" 
  SNSARNUsage    = "The ARN of the receiver SNS topic (mandatory)"
)

flag.StringVar(&arn, "arn", defaultARN, arnUsage)
flag.StringVar(&externalID, "extid", defaultExtID, extIDUsage)
flag.StringVar(&region, "region", defaultRegion, regionUsage)
flag.StringVar(&snsarn, "snsarn", defaultSNSARN, SNSARNUsage)
flag.Parse()

  if snsarn == "" || region == "" {
    fmt.Println("Please supply the mandatory parameters, ARN and region", snsarn, region)
    os.Exit(1)
  }
  createResetTicker()
  http.HandleFunc("/",  func(w http.ResponseWriter, r *http.Request) {
    sendMessage(w, r, snsarn, arn, externalID, region)
  })

  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
 }
}
