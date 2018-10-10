# Overview

Perform requests with filters based on https://www.imdb.com/search/:

## Steps

### Get genders

Perform [Get] on https://www.imdb.com/search/title and filter all types of gender

### For list with n genders get all titles from this until limit end

2 - https://www.imdb.com/search/title?genres={genrer}&adult={include adult results}&sort={sort by}&count={titles per page}

### Create json line with content of search

results/*.jsonl