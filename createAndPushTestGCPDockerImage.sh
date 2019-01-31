cp ~/serviceaccounts/cbh-event-pipeline-16042fb63701.json ./
cp Dockerfiles/GCP-Dockerfile Dockerfile
docker build -t productpromisedeventms .
~/google-cloud-sdk/bin/gcloud docker -- tag productpromisedeventms gcr.io/cbh-event-pipeline/productpromisedeventms:0.0.0.1
~/google-cloud-sdk/bin/gcloud docker -- push gcr.io/cbh-event-pipeline/productpromisedeventms:0.0.0.1
rm ./cbh-event-pipeline-16042fb63701.json
