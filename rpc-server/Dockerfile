FROM golang:1.18
WORKDIR /app
COPY . .
RUN go mod download     # Install Go module dependencies
RUN chmod +x build.sh   # Add execute permission to build.sh
RUN ./build.sh          # Execute build.sh
EXPOSE 8888
CMD ["./output/bootstrap.sh"]
