FROM alpine:latest AS build
RUN apk update
RUN apk upgrade
RUN apk add --update go gcc g++ git
RUN mkdir /app 
ADD snsgateway.go /app/ 
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/service/sns
WORKDIR /app 
#RUN go build -o main .
RUN CGO_ENABLED=1 GOOS=linux go build -o main

FROM alpine:latest
RUN apk add curl
WORKDIR /app
COPY --from=build /app/main .
RUN chmod +x /app/main
ARG REGION
ARG SNSARN
ARG ARN=""
ARG EXTID=""
ARG MAXMESSAGESPERMINUTE=20
ARG PORT=8080
ENV REGION=$REGION
ENV SNSARN=$SNSARN
ENV ARN=$ARN
ENV EXTID=$EXTID
ENV MAXMESSAGESPERMINUTE=$MAXMESSAGESPERMINUTE
ENV PORT=$PORT
ENTRYPOINT /app/main --snsarn "$SNSARN" --region "$REGION" --arn "$ARN" --extid "$EXTID" --maxMessagesPerMinute "$MAXMESSAGESPERMINUTE" --port $PORT
