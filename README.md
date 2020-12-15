## dnspressure

---
DNS pressure test tool for golang.

### Installation
```
go get github.com/devincd/dnspressure
```

### Usage
The dns-pressure has the following features
- Support multiple goroutines pressure measurements
- Support for setting timeout
- Rich results analysis

```
$ dnspressure --help
Usage of dnspressure:
  -concurrent int
    	concurrent users (default 10)
  -domain string
    	The target domain name that needs to be resolved
  -over-time duration
    	over time (default 5s)
  -resp int
    	RESP, number of times to run test in every user (default 500)
  -time duration
    	The total time it takes to execute the test (default 1h0m0s)
```

### Example
```
$ dnspressure --domain=www.baidu.com --over-time=10ms
......
2020/12/15 10:34:57 number: 4995 cost time is 818.25µs
2020/12/15 10:34:57 number: 4996 cost time is 849.394µs
2020/12/15 10:34:57 number: 4997 cost time is 904.567µs
2020/12/15 10:34:57 number: 4998 cost time is 861.168µs
2020/12/15 10:34:57 number: 4999 cost time is 820.395µs
2020/12/15 10:34:57 number: 5000 cost time is 807.793µs
2020/12/15 10:34:57 ----------------- Result Analyse-----------------
2020/12/15 10:34:57 total count		:		 5000
2020/12/15 10:34:57 over time(2ms)  	:		 50
2020/12/15 10:34:57 max time   		:		 7.197964ms
2020/12/15 10:34:57 min time   		:		 0s
2020/12/15 10:34:57 avg time   		:		 0.86ms
2020/12/15 10:34:57 qps        		:		 1159.24q/s
```