# GoFilter
A tool to filter URLs by parameter count or size.


This tool requires unique sorted URL list.  

For example:
```bash
cat hosts.txt | sort -u > sorted && gofilter -f sorted
```

```bash
cat hosts.txt | waybackurls | sort -u > sorted && gofilter -f sorted
```

Usage:
```bash
➜  gofilter ./gofilter -h
Filter given URLs based on query size

Options:
  -f,  --file         File to process
  -c,  --count        Run analysis based on query count, otherwise query length
  -mc, --min-count    Minimum parameter count
  -q,  --quiet        Only print fullpath
```
