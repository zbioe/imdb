# Overview

Perform requests with filters based on https://www.imdb.com/search/:

## Steps

### Get cache and genders

Performs [Get] on https://www.imdb.com/search/title and filter all types of gender

### For list with n genders get all titles from this until limit end

2 - https://www.imdb.com/search/title?genres=film_noir&adult=include&sort=user_rating,desc&count=100

### 