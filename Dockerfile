# Se baseis na imagem https://hub.docker.com/_/golang/
FROM golang:1.8

# Copia o diretorio local para o diretorio do container
ADD . $GOPATH/src/github.com/michelaquino/golang_kafka_example

# Instala a aplicacao
RUN go install github.com/michelaquino/golang_kafka_example

# Executa a aplicacao quando o container for iniciado
ENTRYPOINT $GOPATH/bin/golang_kafka_example

# Expoe a porta 8080
EXPOSE 8080