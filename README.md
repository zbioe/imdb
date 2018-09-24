# IMDB Titles

Simple fetcher for get target tittles with basic infos from [IMDb]("http://www.imdb.com/")

## Requirements

- docker
- docker-compose

## Run
By default will get the 500 most rating titles from each genre:
```sh
make run
```

### Run with options
Options for pass as var for the fetcher change default action.

### Filter by genres (case insensitive)
```sh
make run genre="all"
```

```sh
make run genre="crime,documentary,game_show,mystery,sci_fi"
```

you can see availibles genres

### Change limit
```sh
make run limit="500"
```

## Output

Output will be in [jsonlines](http://jsonlines.org)
and writed inside a results folder in this path

You could use:
```sh
ls results/*.jsonl
````
for see results

