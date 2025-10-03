FROM docker.io/library/alpine:latest
RUN apk --no-cache add ca-certificates gcompat; mkdir /app

# Copy the Pre-built binary file from the previous stage
COPY webserver /app/webserver

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/webserver"]
