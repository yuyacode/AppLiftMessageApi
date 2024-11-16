FROM golang:1.23.3

WORKDIR /app

# airをインストールしてホットリロード対応
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# airでファイル変更を監視しながらアプリケーションを起動
CMD ["air"]

EXPOSE 8080
