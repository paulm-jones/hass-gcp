# hass-gcp

## What is it?
A GCP Cloud Function to consume a [Home Assistant](https://www.home-assistant.io/) message (that has been published to [Pub/Sub](https://cloud.google.com/pubsub/) with the corresponding [Home Assistant plugin](https://www.home-assistant.io/integrations/google_pubsub)), and then insert the contents into BigQuery.
