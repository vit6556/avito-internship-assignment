#!/bin/bash
set -e

/migrator up

exec /http-server