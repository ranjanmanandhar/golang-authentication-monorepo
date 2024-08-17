#!/bin/sh

echo "Sleep For While"
sleep 10

echo "Starting App..."
/usr/local/bin/nettv-auth-consumer $@

