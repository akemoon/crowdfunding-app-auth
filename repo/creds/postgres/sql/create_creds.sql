insert into credentials (
    email,
    password_hash
) values ($1, $2)
returning user_id;
