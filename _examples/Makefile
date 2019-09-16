docker-net:
	docker network create --subnet=172.24.0.0/16 netdemo

docker-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
	docker build -t demo -f Dockerfile.scratch .

docker-start:
ifndef ADDR
	${error 'No address. Usage: `make docker-start ADDR="172.24.0.4" PORT="19000"'}
endif
ifndef PORT
	${error 'No port. Usage: `make docker-start ADDR="172.24.0.4" PORT="19000"'}
endif
	docker run --net netdemo -it -p ${PORT}:${PORT} --ip ${ADDR} -e ADDR=${ADDR} -e PORT=${PORT} demo
