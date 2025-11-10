#!/bin/bash
while IFS=, read -r name code; do
    # do something... Don't forget to skip the header line!
    echo "$code"
done </Users/me/Developer/go-projects/news-sites/sql/domains_to_county.txt
