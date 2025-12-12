# Stage 1: Builder for Go application
FROM golang:1.24-alpine AS builder

ENV HTTP_PROXY=http://192.168.111.5:8080
ENV HTTPS_PROXY=http://192.168.111.5:8080
ENV NO_PROXY=127.0.0.1,localhost,70.170.65.133

WORKDIR /app

# Install build dependencies for Go
#RUN apk add --no-cache --no-check-certificate git gcc libc-dev
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update --no-check-certificate && apk add --no-cache --no-check-certificate ca-certificates git gcc libc-dev

# Copy go.mod and go.sum first to leverage Docker cache
COPY backend/go.mod backend/go.sum ./
RUN go env -w GO111MODULE=on &&  go env -w GOPROXY=https://mirrors.xfusion.com/goproxy/ && go env -w GONOSUMDB=* && go mod download
#RUN GOPROXY=https://goproxy.cn,direct \
#    GOSUMDB=off \
#    GOINSECURE=* \
#    GIT_SSL_NO_VERIFY=true \
#    GOFLAGS="-modcacherw" \
#    go env -w GONOSUMDB="*" && \
#    go mod download

# Copy the entire backend source code
COPY backend/ ./

# Build the Go application
RUN go build -ldflags="-s -w" -o /app/opencoze main.go


# Stage 2: Final image
FROM alpine:3.22.0

ENV HTTP_PROXY=http://192.168.111.5:8080
ENV HTTPS_PROXY=http://192.168.111.5:8080
ENV NO_PROXY=127.0.0.1,localhost,70.170.65.133

WORKDIR /app

# Install runtime dependencies for Go app and base for Python
# pax-utils for scanelf, python3 for running Python, python3-dev for headers/shared libs
# bind-tools for nslookup etc., file for debugging file types
#RUN apk add --no-cache --no-check-certificate ca-certificates pax-utils python3 python3-dev bind-tools file deno curl
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update --no-check-certificate && apk add --no-cache --no-check-certificate ca-certificates pax-utils python3 python3-dev bind-tools file deno curl

# RUN deno run -A jsr:@langchain/pyodide-sandbox -c "print('Hello, World')"

# Install Python build dependencies, create venv, install packages, then remove build deps
#RUN apk add --no-cache --no-check-certificate --virtual .python-build-deps build-base python3-dev py3-pip git && \
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update --no-check-certificate && apk add --no-cache --no-check-certificate --virtual .python-build-deps build-base python3-dev py3-pip git && \
    python3 -m venv --copies --upgrade-deps /app/.venv && \
    # Activate venv and install packages
    . /app/.venv/bin/activate && \
    # If you want to use other third-party libraries, you can install them here.
    pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn urllib3==1.26.16 && \
    pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn --no-cache-dir h11==0.16.0 httpx==0.28.1 pillow==11.2.1 pdfplumber==0.11.7 python-docx==1.2.0 numpy==2.3.1 && \
    # Deactivate (optional, as RUN is a new shell)
    # deactivate && \
    # Remove build dependencies
    apk del --no-check-certificate .python-build-deps


# Copy the built Go binary from the builder stage
COPY --from=builder /app/opencoze /app/opencoze

# Copy Python application scripts
COPY backend/infra/document/parser/impl/builtin/parse_pdf.py /app/parse_pdf.py
COPY backend/infra/document/parser/impl/builtin/parse_docx.py /app/parse_docx.py
COPY backend/infra/coderunner/impl/script/sandbox.py /app/sandbox.py


# Copy static resources
# COPY backend/static /app/resources/static/
COPY backend/conf /app/resources/conf/
COPY docker/.env.example /app/.env
# COPY docker/.env.ve /app/.env
# COPY docker/cert.pem /app/cert.pem
# COPY docker/key.pem /app/key.pem

# Set PATH to prioritize venv's binaries
ENV PATH="/app/.venv/bin:${PATH}"
# ENV LD_LIBRARY_PATH="/usr/lib:${LD_LIBRARY_PATH}" # Keep commented for now

# Ensure python scripts and venv executables are executable
RUN chmod +x /app/parse_pdf.py /app/parse_docx.py  && \
    find /app/.venv/bin -type f -exec chmod +x {} \;


EXPOSE 8888

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

ENV HTTP_PROXY=""
ENV HTTPS_PROXY=""
ENV NO_PROXY=""

CMD ["/app/opencoze"]
