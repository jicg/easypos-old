# The FROM instruction sets the Base Image for subsequent instructions.
# Using Nginx as Base Image
FROM golang:1.7.1
MAINTAINER <284077318@qq.com>
COPY . $GOPATH/src/easypos
# The RUN instruction will execute any commands
# Adding HelloWorld page into Nginx server
WORKDIR $GOPATH/src/easypos
RUN go get easypos
RUN go install -a easypos

# The EXPOSE instruction informs Docker that the container listens on the specified network ports at runtime
EXPOSE 4000

# The CMD instruction provides default execution command for an container
# Start Nginx and keep it from running background
CMD easypos web --port 4000