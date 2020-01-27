DROP TABLE services;
CREATE TABLE services
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       varchar(255)  not null,
    url        varchar(2083) not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    updated_at datetime default CURRENT_TIMESTAMP not null,
    deleted_at datetime      null
);

DROP TABLE service_ticks;
CREATE TABLE service_ticks
(
    service_id       integer not null references services (id),
    is_up            bool    not null,
    response_time_ms float64 not null,
    created_at       datetime default CURRENT_TIMESTAMP not null
);

DROP TABLE service_metrics;
CREATE TABLE service_metrics
(
    service_id               integer not null references services (id) primary key,
    tick_count               integer not null default 0,
    up_count                 integer not null default 0,
    down_count               integer not null default 0,
    average_response_time_ms double  not null default 0,
    created_at               datetime         default CURRENT_TIMESTAMP not null,
    updated_at               datetime         default CURRENT_TIMESTAMP not null
);

DROP TABLE service_events;
CREATE TABLE service_events
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id   integer     not null references services (id),
    event        varchar(20) not null,
    date_started datetime    not null default CURRENT_TIMESTAMP,
    date_ended   datetime    null,
    created_at   datetime             default CURRENT_TIMESTAMP not null,
    updated_at   datetime             default CURRENT_TIMESTAMP not null
);

INSERT INTO services (id, name, url)
VALUES (NULL, "email-service", "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/email-service/health"),
       (NULL, "payment-gateway", "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/payment-gateway/health"),
       (NULL, "microservice-controller",
        "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/microservice-controller/health"),
       (NULL, "transaction-monitor",
        "https://u0e8utqkk2.execute-api.eu-west-2.amazonaws.com/dev/transaction-monitor/health");

