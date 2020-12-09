#!/usr/bin/env bash

echo "====================="
echo "=== End2End tests ==="
echo "====================="
echo "python3 -m unittest test_alone.py -v"
python3 -m unittest test_alone.py -v

echo "python3 -m unittest test_winner.py -v"
python3 -m unittest test_winner.py -v