#!/bin/bash

mkdir -p ./static

# HTMX
wget https://unpkg.com/htmx.org/dist/htmx.min.js -O ./static/htmx.min.js

# Alpine.js
wget https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js -O ./static/alpine.min.js

# Preline UI
wget https://raw.githubusercontent.com/htmlstreamofficial/preline/main/dist/dropdown.js -O ./static/dropdown.js
