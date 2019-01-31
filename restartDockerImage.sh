cp ~/serviceaccounts/cbh-event-pipeline-16042fb63701.json ./
cp Dockerfiles/test-Dockerfile Dockerfile
docker build -t productpromisedeventms .
docker stop productpromisedeventms
docker rm productpromisedeventms
docker run -d --restart=always --publish 8093:8080 --name productpromisedeventms productpromisedeventms
rm ./cbh-event-pipeline-16042fb63701.json