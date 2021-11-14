FROM golang:1.17 as build

WORKDIR /app

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/prisma/prisma-client-go prefetch

COPY . ./

# create and migrate the database
RUN go run github.com/prisma/prisma-client-go migrate

# generate the Prisma Client Go client
RUN go run github.com/prisma/prisma-client-go generate
# or, if you use go generate to run the generator, use the following line instead
# RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /main .

ARG PORT
ARG DATABASE_URL
ARG DISCORD_TOKEN
ARG LINE_CHANNEL_SECRET
ARG LINE_CHANNEL_TOKEN
ARG GUILD_ID
ARG PARENT_ID

ENV PORT $PORT
ENV DATABASE_URL $DATABASE_URL
ENV DISCORD_TOKEN $DISCORD_TOKEN
ENV LINE_CHANNEL_SECRET $LINE_CHANNEL_SECRET
ENV LINE_CHANNEL_TOKEN $LINE_CHANNEL_TOKEN
ENV GUILD_ID $GUILD_ID
ENV PARENT_ID $PARENT_ID

CMD ["/main"]
