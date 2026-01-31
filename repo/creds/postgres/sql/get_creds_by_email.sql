select user_id,
       email,
       password_hash
from credentials
where email = $1
