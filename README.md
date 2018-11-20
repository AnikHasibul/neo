# neo

> Unfortunately, no one can be told what The Matrix is. You have to see it for yourself. -neo (The Matrix 1999)

## Installation

Install via `go get`:

```txt
go get -v github.com/anikhasibul/neo
```

Verify the installation:

```txt
$ neo --help
```

## Let `neo` say

Now open a file to let neo print the file to read it yourself!

```txt
$ neo longlonglongfile.txt
```

The default maxspeed is 100ms,, you can set it yourself.

```txt
$ neo -s 200 longlonglongfile.txt

```

The default color is green. You can use other color or no color:

Green color:

```txt
$ neo -c 1 longlonglongfile.txt
```

Red color:

```txt
$ neo -c 2 longlonglongfile.txt
```

for more colors see `neo --help`

## How to quite neo!

* To pause `neo` presss `CTRL`+`C`.

* To resume the output press any key + ENTER.

* To quite `neo` press `CTRL`+`C`+`q` and enter.
