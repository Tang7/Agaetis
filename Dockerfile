FROM hopipola/tf_go

WORKDIR /go/agaetis

COPY . .

RUN mkdir -p /go/agaetis/model && \
  curl -o /go/agaetis/model/inception5h.zip -s "http://download.tensorflow.org/models/inception5h.zip" && \
  unzip /go/agaetis/model/inception5h.zip -d /go/agaetis/model/inception5h

RUN go build

EXPOSE 8080

ENTRYPOINT ["./agaetis"]