FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y

ENV GOPATH /home/cfrye/go
WORKDIR /home/cfrye/go/src/github.com/cfrye2000/productPromisedEventMS

ADD . $GOPATH/src/github.com/cfrye2000/productPromisedEventMS

RUN cd $GOPATH/src/github.com/cfrye2000/productPromisedEventMS

EXPOSE 8080

RUN cp $GOPATH/src/github.com/cfrye2000/productPromisedEventMS/configs/gcpProductPromisedEventMS.cfg $GOPATH/src/github.com/cfrye2000/productPromisedEventMS/productPromisedEventMS.cfg

ADD cbh-event-pipeline-16042fb63701.json $GOPATH/src/github.com/cfrye2000/productPromisedEventMS/

RUN go get -u cloud.google.com/go/pubsub

ENV GOOGLE_APPLICATION_CREDENTIALS $GOPATH/src/github.com/cfrye2000/productPromisedEventMS/cbh-event-pipeline-16042fb63701.json

RUN go install

# Run the outyet command by default when the container starts.
CMD $GOPATH/bin/productPromisedEventMS