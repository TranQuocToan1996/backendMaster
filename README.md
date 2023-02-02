

docker build -t simplebank:latest .
docker run --name simplebank -p 8080:8080 simplebank
