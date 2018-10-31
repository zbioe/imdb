# IMDB Titles

Simple fetcher for get target tittles with basic infos from [IMDb]("http://www.imdb.com/")

## Requirements

docker or go

## Build
For get imdb binary in folder, you can use:
```sh
make build
```

## Test

Run all tests
```sh
make check
```

### Test with args
for pass args to test
```sh
make check-args args="-v genrer/*.go -run TestParse"
# or
make check-args args="./... -bench Parse"
```

## Run
By default will get the 500 most rating titles from each genre:

```sh
make run
```

### Run with options
Options for pass as var for the fetcher change default action.
```
imdb -h
```

for pass args to makefile you can set flags in var args:
```
make run args='--limit=1000 --sort="num_votes,asc" --debug --adult=false'
```

### Change limit
for change max limit, you can use flag limit
```sh
imdb --limit=1000
```

### Remove adult results
for prevent adult titles you can use flag adult
```sh
imdb --adult=false
```

### Change sort
if you need change sort, you can use flag sort
```sh
imdb --sort="num_votes,asc"
```

### Debug mode
for activate debug mode, you can use flag debug
```sh
imdb --debug
```

## Output

Output will be in [jsonlines](http://jsonlines.org)
and writed inside a results folder in current path

for see results, you can use:
```sh
ls results/*.jsonl
````

### Example 

```json
{
  "Name": "Apenas Um Show",
  "Episode": "A Regular Epic Final Battle",
  "Year": "(2009–2017)",
  "Genres": [
    "animation",
    "action",
    "adventure"
  ],
  "Rating": {
    "Value": 9.9,
    "Best": 10,
    "Count": 571,
    "Position": 1
  }
}
```

