# GoFilter
A tool to filter URLs by parameter count or size.


### Examples:
```bash
cat hosts.txt | GoFilter -s
```

```bash
cat hosts.txt | waybackurls | GoFilter -s
```

```bash
GoFilter -f list.txt
```

### Installation: 
```bash
export GO111MODULE=off; go get github.com/cryonayes/GoFilter
```

### Usage:
```bash
➜  GoFilter -h
Filter given URLs based on query size

Options:
  -f,  --file          File to process
  -s   --std           Read from standard input
  -l,  --length        Run analysis based on query length
  -mc, --min-count     Minimum parameter count
  -q,  --quiet         Only print fullpath
```
