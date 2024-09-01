# ベースイメージ
FROM golang:1.23.0

# 作業ディレクトリの設定
WORKDIR /app

# go.modとgo.sumをコピーして依存関係を解決
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# airをインストール (新しいリポジトリを使用)
RUN go install github.com/air-verse/air@latest

# 環境変数PATHを設定
ENV PATH="/go/bin:${PATH}"

# ホットリロードのためにairを実行
CMD ["air"]