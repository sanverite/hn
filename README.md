# hn

A fast Hacker News client for your terminal. Fetch top stories,
Ask HN, and open links in your browser without leaving the command line.

## Installation

```bash
git clone https://github.com/sanverite/hn
cd hn
make build
./bin/hn --help
```

## Usage

```bash
# Fetch top 10 stories
./bin/hn top

# Fetch top 30 stories
./bin/hn top --limit 30

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
