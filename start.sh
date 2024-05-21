#!/usr/bin/env bash

Xvfb -ac :99 -screen 0 1280x1024x16 &
DISPLAY=:99 exec ./app :8080