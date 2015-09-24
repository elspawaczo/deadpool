CREATE TABLE report (
    Id                  serial primary key,
    origin              text,
    method              text,
    status              int,
    content_type        text,
    content_length      bigint,
    host                text,
    url                 text,
    scheme              text,
    path                text,
    body                text,
    request_body        text,
    date_start          timestamp,
    date_end            timestamp,
    time_taken          timestamp
);
