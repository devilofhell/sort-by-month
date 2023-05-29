:: build and run a container on your local machine
:: if you are executing this script in a windows environment, the local go-language needs to compile it for linux
go env -w GOOS=linux
go build
go env -w GOOS=windows
:: docker compose up -d --build
docker rm -f sort-by-month-c
docker build -t 559rvsuq/sort-by-month:latest .
docker run -d -i -t --name "sort-by-month-c" 559rvsuq/sort-by-month:latest