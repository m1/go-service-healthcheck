# go-service-healthcheck

# How it works

## Runner
Basically the runner works by querying the service every defined time in the .env. Collecting both
if the service is up and how long it took to respond. 

It will then update the `service_metric` - updating the up count and calculating the rolling average response
time. 

It will then update the `service_events` - if the service is the same "event" i.e. up or down as the last
event it will do nothing. However if the status of the service has changed i.e. from up to down, it will
then update the old event ending the "uptime" and adding a new downtime event.

It will then finally add a `service_tick" which basically just saves the raw data for each tick - just the
up/down status and the response time.

## API
The API just pulls the data from above and does a few parsing operations, i.e. marshalling all the times
to unix timestamps. 

For the `/v1/services` endpoint it just displays all the scrape service targets, their metrics and the latest
event. 

For the `/v1/services/{id}/events` endpoint it just pulls all the events for the corresponding service id. 

# How to run

1. `make local-run` or `make docker-up` or just `go build && ./go-service-healthcheck run` 
2. `curl -X GET http://localhost:9999/v1/services`

# Config

Create an `.env` with:
```.env
DEBUG=true
ENV=development
SERVICE_NAME=go-service-healthcheck

API_PORT=9999	router.Get("/{serviceID}/metrics", h.GetServices)

API_DOMAIN=localhost

DB_FILE=db.db

RUNNER_SCRAPE_INTERVAL_SECONDS=15
RUNNER_HTTP_TIMEOUT_SECONDS=15
```

# Endpoints 

## /v1/services 

Example query:
```
curl -X GET http://localhost:9999/v1/services
```

Example response:
```json
{
    "status": 200,
    "status_desc": "OK",
    "data": {
        "services": {
            "data": [
                {
                    "created_at": 1579723302,
                    "current_status": {
                        "created_at": 1579723308,
                        "date_ended": null,
                        "date_started": 1579723308,
                        "event": "uptime",
                        "id": 2,
                        "updated_at": 1579723308
                    },
                    "id": 1,
                    "latest_event": {
                        "data": {
                            "created_at": 1579723308,
                            "date_ended": null,
                            "date_started": 1579723308,
                            "event": "uptime",
                            "id": 2,
                            "updated_at": 1579723308
                        }
                    },
                    "metric": {
                        "data": {
                            "average_response_time_ms": 54.52497556578949,
                            "down_count": 0,
                            "tick_count": 75,
                            "up_count": 75,
                            "uptime_percent": 100
                        }
                    },
                    "name": "email-service",
                    "updated_at": 1579723302,
                    "url": "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/email-service/health"
                },
                {
                    "created_at": 1579723302,
                    "current_status": {
                        "created_at": 1579723308,
                        "date_ended": null,
                        "date_started": 1579723308,
                        "event": "downtime",
                        "id": 3,
                        "updated_at": 1579723308
                    },
                    "id": 2,
                    "latest_event": {
                        "data": {
                            "created_at": 1579723308,
                            "date_ended": null,
                            "date_started": 1579723308,
                            "event": "downtime",
                            "id": 3,
                            "updated_at": 1579723308
                        }
                    },
                    "metric": {
                        "data": {
                            "average_response_time_ms": 53.46076003896105,
                            "down_count": 76,
                            "tick_count": 76,
                            "up_count": 0,
                            "uptime_percent": 0
                        }
                    },
                    "name": "payment-gateway",
                    "updated_at": 1579723302,
                    "url": "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/payment-gateway/health"
                },
                {
                    "created_at": 1579723302,
                    "current_status": {
                        "created_at": 1579723308,
                        "date_ended": null,
                        "date_started": 1579723308,
                        "event": "downtime",
                        "id": 1,
                        "updated_at": 1579723308
                    },
                    "id": 3,
                    "latest_event": {
                        "data": {
                            "created_at": 1579723308,
                            "date_ended": null,
                            "date_started": 1579723308,
                            "event": "downtime",
                            "id": 1,
                            "updated_at": 1579723308
                        }
                    },
                    "metric": {
                        "data": {
                            "average_response_time_ms": 33.92872929870131,
                            "down_count": 76,
                            "tick_count": 76,
                            "up_count": 0,
                            "uptime_percent": 0
                        }
                    },
                    "name": "microservice-controller",
                    "updated_at": 1579723302,
                    "url": "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/microservice-controller/health"
                },
                {
                    "created_at": 1579723302,
                    "current_status": {
                        "created_at": 1579723793,
                        "date_ended": null,
                        "date_started": 1579723793,
                        "event": "uptime",
                        "id": 5,
                        "updated_at": 1579723793
                    },
                    "id": 4,
                    "latest_event": {
                        "data": {
                            "created_at": 1579723793,
                            "date_ended": null,
                            "date_started": 1579723793,
                            "event": "uptime",
                            "id": 5,
                            "updated_at": 1579723793
                        }
                    },
                    "metric": {
                        "data": {
                            "average_response_time_ms": 2470.3923702837833,
                            "down_count": 28,
                            "tick_count": 73,
                            "up_count": 45,
                            "uptime_percent": 61.64383561643836
                        }
                    },
                    "name": "transaction-monitor",
                    "updated_at": 1579723302,
                    "url": "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/transaction-monitor/health"
                }
            ]
        }
    }
}
```

Few notes on response:
- `latest_event`: the runner saves 'events' i.e. from going from 'down' to 'up'
- `metric`: This is just the data structure that holds the accumulated data about the service, i.e. avg response time,
what percent it has been up etc.

## /v1/services/{id}/events

Example query:
```
curl -X GET http://localhost:9999/v1/services/4/events
```

Example response:
```json
{
    "status": 200,
    "status_desc": "OK",
    "data": {
        "events": {
            "data": [
                {
                    "created_at": 1579723793,
                    "date_ended": null,
                    "date_started": 1579723793,
                    "event": "uptime",
                    "id": 5,
                    "updated_at": 1579723793
                },
                {
                    "created_at": 1579723314,
                    "date_ended": 1579723793,
                    "date_started": 1579723314,
                    "event": "downtime",
                    "id": 4,
                    "updated_at": 1579723793
                }
            ]
        }
    }
}
```

Here we can see we had downtime between `1579723793` and `1579723314` but then came back up.

# TODOs
- More tests (primarily runner.go and worker.go)!
- Replace sqlite with proper db - time series db would be a good fit (influxdb?)
- Abstract out service repo to distinct repos i.e. ServiceMetricRepo
- Make error checking/catching a bit better - i.e catch when services are erroring
- Move to worker pool implementation for parsing services
- Split runner and api into two services
- Add querying of services (i.e. prometheus/promql)
- Paginate the events and services