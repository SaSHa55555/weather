FROM golang
WORKDIR /doc_weather
COPY . /doc_weather/
EXPOSE 3000
CMD [ "go", "run", "main.go" ]