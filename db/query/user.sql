-- name: ListUser :many
SELECT * FROM kuser;

-- name: CreateUser :exec
INSERT INTO kuser (FirstName, LastName, Phone, Addr)
VALUES ($1, $2, $3, $4);

-- name: DeactivateUser :exec
UPDATE kuser
SET Active = False
WHERE Phone = $1;

-- name: DeleteUser :exec
delete from kuser where phone = '7408963464';