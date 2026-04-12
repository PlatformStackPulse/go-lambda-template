FROM golang:1.23 AS builder

ARG TARGETARCH=arm64
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
ARG GO_VERSION=unknown

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X github.com/PlatformStackPulse/go-lambda-template/pkg/version.Version=${VERSION} -X github.com/PlatformStackPulse/go-lambda-template/pkg/version.Commit=${COMMIT} -X github.com/PlatformStackPulse/go-lambda-template/pkg/version.BuildTime=${BUILD_TIME} -X github.com/PlatformStackPulse/go-lambda-template/pkg/version.GoVersion=${GO_VERSION}" \
    -o /out/bootstrap \
    ./cmd/lambda

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /out/bootstrap ${LAMBDA_TASK_ROOT}/bootstrap
