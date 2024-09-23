-- Write your migrate up statements here

CREATE TABLE messages(
    id varchar(255),
    negotiation_id varchar(255) NOT NULL REFERENCES negotiations(id),
    sender_email varchar(255) NOT NULL,
    sender_name varchar(255) NOT NULL,
    message_text varchar(255) NOT NULL,
    time_sent TIMESTAMP DEFAULT NOW(),
    status varchar(255) DEFAULT 'Sent',
    PRIMARY KEY(id)
);

---- create above / drop below ----

DROP TABLE messages;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
