# FROM    golang:1.19.3 as base
# WORKDIR /app

# COPY    go.mod ./
# COPY    go.sum ./
# RUN     go mod download

# COPY    * ./
# COPY    *.env ./

# # RUN     go build -o last.go

# EXPOSE  8080

# RUN go build -v -o apka

# CMD ["apka"]

# CMD     ["./last.go"]


FROM    golang:1.19.3 as base
WORKDIR /app

COPY    . ./
RUN     go mod download

RUN     go build -o app

EXPOSE  8080

CMD     ["./app"]