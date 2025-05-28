FROM alpine:latest

RUN mkdir /app 

# copy the binary from the builder stage to the new image
COPY loggerServiceApp /app

CMD ["/app/loggerServiceApp"]