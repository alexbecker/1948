#! /bin/bash
source env.sh
sqlite3 $DATABASE < schema.sql
