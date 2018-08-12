#### SNS Gateway

This solution is created to allow applications without SNS support, to send messages with webhooks, using either parameters or a simple request body. The solution can use IAM roles to authenticate to AWS and it can limit the number of messages sent per minute.

### Quick start for testing
Clone the repo and run the snsgateway.go. Example:

`go run snsgateway.go --snsarn arn:aws:sns:yourregion:youraccountnumber:topicname --region yourregion`

### Prerequisites
```
  yum install go
  go get github.com/aws/aws-sdk-go/aws
  go get github.com/aws/aws-sdk-go/service/s3
  go get github.com/aws/aws-sdk-go/service/sns
```

### Options

```
  -arn string
        The ARN of the role you want to assume
  -extid string
        The ExternalID constraint, if applicable for the role you need to assume
  -maxMessagesPerMinute int
        The maximum number of messages allowed per minute (default 20)
  -port int
        The listening port for the application (default 8080)
  -region string
        The region of the SNS topic (mandatory)
  -snsarn string
        The ARN of the receiver SNS topic (mandatory)
```

### Usage
```
Send a JSON Message in the request body with default Subject (FromSNSGateway)
curl 127.0.0.1:8080 -d '{test: "test"}'

Send a Message by using parameter with default Subject
curl 127.0.0.1:8080?message=test

Send a Message by using parameters with Subject "test"
curl '127.0.0.1:8080?message=test&subject=test'

Send a JSON Message in the request body with Subject "test"
curl '127.0.0.1:8080?subject=test' -d '{test: "test"}'
```

### Deployment

The solution can be run on Kubernetes.
1. Create a namespace (or use an existing one)
2. Create a configmap to provide the ENV Variables
3. Create the deployment and the service to expose the port inside you k8s cluster

```
kubectl create namespace monitoring
kubectl -f create 1_env.yaml -n monitoring
kubectl -f create 2_snsgateway.yaml -n monitoring
```

### Security
The solution is using standard IAM based security to provide access to SNS. If the container has an Instance Profile it will be automatically picked up. Otherwise you have to create a .aws/credentials file on the container, or to provide a Role ARN for the Assumerole with the arn parameter. You can read about this in more detail here: [link](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

The solution does not provide any SSL termination or authorization/authentication for the listening endpoint. You can set up a reverse proxy to secure the endpoint. Example: [ingress-nginx](https://github.com/kubernetes/ingress-nginx) 

### Acknowledgments

    - Nick Gauthier (https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/)
    - William Kennedy (https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html)
    - Edd Turtle (https://golangcode.com/get-a-url-parameter-from-a-request/)
    - Adam Crosby (https://github.com/adamcrosby/sts-example)
    
### License

This project is licensed under the MIT License - see the LICENSE.md file for details
