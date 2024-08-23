# Lightspeed exercise

# Unique IP Address Counter 

You have a simple text file with IPv4 addresses. One line is one address, line by line:

```
145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5
...
```

The file is unlimited in size and can occupy tens and hundreds of gigabytes.

You should calculate the number of __unique addresses__ in this file using as little memory and time as possible.
There is a "naive" algorithm for solving this problem (read line by line, put lines into HashSet).
It's better if your implementation is more complicated and faster than this naive algorithm.


---
Before submitting an assignment, it will be nice to check how it handles this [file](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip). Attention - the file weighs about 20Gb, and unzips to about 120Gb.


# Solution

The program supports two algorithms for counting unique addresses:

1.	[Bloom Filter](https://en.wikipedia.org/wiki/Bloom_filter) - A probabilistic data structure that efficiently checks whether an IP address has already been counted. This method is memory-efficient but may produce false positives, leading to a slight inaccuracy in the count.

2.	Simple "naive" algorithm with map - A deterministic algorithm that uses a thread-safe sync.Map structure to store all unique IP addresses. This method is accurate but may require more memory and time, especially for large files.

# How to run

The program accepts two command-line parameters:

```
1) The name of the file containing IP addresses.

2) The algorithm to use for counting:
   "1" - for using the Bloom Filter (default).
   "2" - for using the "naive" algorithm with map.
```
Example:
```bash
bin/counter result.txt 1
```

# Results

On MacOS Apple M1 Max 32 GB memory with NumCPU=10 the results are as follows: 

| File size (MB) | "naive" algorithm |         | Bloom Filter |       |
|----------------|-------------------|---------|--------------|-------|
|                | result            | time    | result       | time  |
|                |                   |         |              |       |
| 10             | 734 302           | 643ms   | 734 302      | 420ms |
| 100            | 7 336 499         | 7.79s   | 7 336 464    | 3.33s |
| 1 000          | 72 798 930        | 2m24s   | 72 755 437   | 43s   |
| 10 000         | 610 846 107       | more 2h | 610 836 289  | 6m20s |


Generation of file with size of 10 GB with random IP addresses take time about 27m 53s.
