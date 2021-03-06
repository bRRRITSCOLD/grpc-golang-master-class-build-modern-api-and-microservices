#!/usr/bin/env bash

printf "setting up mongo users\n"
mongo localDb <<EOF
db.createRole({
  role: "readWriteMinusDropRole",
  privileges: [
  {
    resource: { db: "localDb", collection: ""},
    actions: [ "collStats", "dbHash", "dbStats", "find", "killCursors", "listIndexes", "listCollections", "convertToCapped", "createCollection", "createIndex", "dropIndex", "insert", "remove", "renameCollectionSameDB", "update"]} ],
    roles: []
  }
);
use admin;
db.createUser({user: "localuser", pwd: "1234abcd", roles: [{role: "readWriteMinusDropRole", db: "localDb"}]})
quit()
EOF
printf "set up mongo users\n"
