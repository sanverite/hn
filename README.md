# hn

A fast Hacker News client for your terminal. Fetch top, best, and new
stories, then open links in your browser without leaving the command line.

## Installation

```bash
git clone https://github.com/sanverite/hn
cd hn
make build
./bin/hn --help
```

## Usage

```bash
# Top stories (default 10)
./bin/hn top

# Top 30 stories
./bin/hn top --limit 30

# Best stories
./bin/hn best

# Best 20 stories
./bin/hn best --limit 20

# New stories
./bin/hn new

# New 15 stories
./bin/hn new --limit 15

# Open a story in your browser by ID
./bin/hn open 43291234

# Print build info
./bin/hn version
```

## Development

```bash
make build   # build the binary to bin/hn
make test    # run all tests
make race    # run tests with the race detector
```
