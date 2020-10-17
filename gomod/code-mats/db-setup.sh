#!/bin/bash
createdb -h localhost -U postgres code-mats
psql -h localhost -U postgres -f create_table.sql code-mats

