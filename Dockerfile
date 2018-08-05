FROM alpine:3.8
ADD main /
RUN apk --no-cache add curl go
RUN chmod +x /main
CMD ["/main"]
