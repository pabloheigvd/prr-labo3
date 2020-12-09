#!/usr/bin/env bash

echo "====================="
echo "=== End2End tests ==="
echo "====================="
echo "python3 -m unittest test_sync.py -v"
python3 -m unittest test_alone.py -v
