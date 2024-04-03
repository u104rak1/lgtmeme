\set pgpass `echo "password"`

ALTER USER supabase_storage_admin WITH PASSWORD :'pgpass';