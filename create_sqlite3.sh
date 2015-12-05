#!/bin/sh

db="new_empty_db.sqlite3"

echo "try to create sqlite3 database"
echo $db

if [ -e $db ]; then
    echo "database already exist"
    exit 0
fi

echo "\
CREATE TABLE teleport ( \
    encrypted_addr int NOT NULL PRIMARY KEY,
    encrypted_private_key varchar(32) NOT NULL
); \
" | sqlite3 $db

echo "success create a new empty schema"
