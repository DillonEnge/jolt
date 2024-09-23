-- name: RecordMessage :one
INSERT INTO messages(id, negotiation_id, sender_email, sender_name, message_text)
VALUES(
    uuid_generate_v4(),
    @negotiation_id::text,
    @sender_email::text,
    @sender_name::text,
    @message_text::text
)
ON CONFLICT(id) DO UPDATE SET
    message_text = excluded.message_text,
    status = excluded.status
RETURNING *;

-- name: MessagesByNegotiationID :many
SELECT m.* FROM messages m
LEFT JOIN negotiations n ON m.negotiation_id = n.id
WHERE n.id = @negotiation_id::text;
