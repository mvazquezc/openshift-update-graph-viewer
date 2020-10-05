FROM golang:latest
RUN go get github.com/mvazquezc/openshift-update-graph-viewer && go get github.com/gorilla/mux
WORKDIR /go/src/github.com/mvazquezc/openshift-update-graph-viewer/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=0 /go/src/github.com/mvazquezc/openshift-update-graph-viewer/main .
COPY --from=0 /go/src/github.com/mvazquezc/openshift-update-graph-viewer/templates .
COPY --from=0 /go/src/github.com/mvazquezc/openshift-update-graph-viewer/templates/index-template.html ./templates/
EXPOSE 8080
CMD ["/main"]
