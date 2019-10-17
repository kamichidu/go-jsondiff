[![Build Status](https://travis-ci.org/kamichidu/go-jsondiff.svg?branch=master)](https://travis-ci.org/kamichidu/go-jsondiff)
[![Coverage Status](https://coveralls.io/repos/github/kamichidu/go-jsondiff/badge.svg)](https://coveralls.io/github/kamichidu/go-jsondiff)
[![godoc](https://godoc.org/github.com/kamichidu/go-jsondiff?status.svg)](https://godoc.org/github.com/kamichidu/go-jsondiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamichidu/go-jsondiff)](https://goreportcard.com/report/github.com/kamichidu/go-jsondiff)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://raw.githubusercontent.com/kamichidu/go-jsondiff/master/LICENSE)

# go-jsondiff
yet another json diff tool for golang

# CLI tool

Install via `go get`

```
go get github.com/kamichidu/go-jsondiff/json-diff
```

Example diff output:

```a.json
{
    "level": "debug",
    "log": {
        "required": ["hoge", "fuga"],
        "properties": {
            "hoge": {
                "type": "string"
            },
            "fuga": {
                "type": "object",
                "required": ["a", "b"],
                "properties": {
                    "a": {},
                    "b": {}
                }
            }
        }
    },
    "title": "Internal server error",
    "bool": true,
    "float": 0.0,
    "int": 0,
    "eles": [{
        "msg": "a"
    }, {
        "msg": "b"
    }]
}
```
```b.json
{
    "level": "info",
    "log": {
        "required": ["fuga", "hoge"],
        "properties": {
            "hoge": {
                "type": "string"
            },
            "fuga": {
                "type": "object",
                "required": ["b", "a"],
                "properties": {
                    "a": {},
                    "b": {}
                }
            }
        }
    },
    "title": "Internal server error\nhello world",
    "bool": false,
    "float": 1.0,
    "int": 1,
    "eles": [{
        "msg": "b"
    }, {
        "msg": "a"
    }, {
        "msg": "c"
    }]
}
```
```
> json-diff --set-property '$..required' --set-property '$.eles' a.json b.json
--- "a.json"    2019-03-01 16:29:21.065481500 +0900
+++ "b.json"    2019-03-01 16:48:58.315785000 +0900
@@ $.bool @@
-true
+false

@@ $.eles @@
+{"msg":"c"}

@@ $.float @@
-0.0
+1.0

@@ $.int @@
-0
+1

@@ $.level @@
-"debug"
+"info"

@@ $.title @@
-"Internal server error"
+"Internal server error\nhello world"
```
