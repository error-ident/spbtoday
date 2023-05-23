build:
	docker build -f ./deployments/Dockerfile -t spbtoday:1.0.0 .

run: build
	docker-compose -f ./deployments/docker-compose.yaml -p spbtoday up -d