# Using golang to publish on GCP PubSub

This is a small snippet to publish an message on GCP Pubsub



##Pre requisite:
- Service account of a GCP project
- A PubSub Topic created in a GCP Project

## Environment variables
- `GOOGLE_APPLICATION_CREDENTIALS`. The service account JSON
- `PROJECT_ID`. GCP's project id
- `TOPIC_ID`. PubSub Topic ID
