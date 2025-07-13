-- name: CreateTicket :exec

INSERT INTO tickets (
    full_name, 
    email_address, 
    issue_category, 
    priority, 
    issue_summary, 
    detailed_description,
    department,
    created_at, 
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    now(), 
    now()
);

-- name: GetAllTickets :many
SELECT * FROM tickets;

-- name: GetTicket :one
SELECT * FROM tickets 
WHERE id = $1;