# EMDR to NSQ
[![Build Status](https://drone.element-43.com/api/badges/EVE-Tools/emdr-to-nsq/status.svg)](https://drone.element-43.com/EVE-Tools/emdr-to-nsq) [![Go Report Card](https://goreportcard.com/badge/github.com/eve-tools/emdr-to-nsq)](https://goreportcard.com/report/github.com/eve-tools/emdr-to-nsq) [![Docker Image](https://images.microbadger.com/badges/image/evetools/emdr-to-nsq.svg)](https://microbadger.com/images/evetools/emdr-to-nsq)

This is a very simple service for [Element43](https://element-43.com) which just connects EMDR to our internal NSQ. It aims to be way faster than solutions we had before. For now it only supports `orders` messages from EMDR. EMDR's messages are split up by rowset before being submitted to the `orders` queue. Also, the order's attributes get mapped for easier access later on. Processing bulk updates of the market can result in lots of messages (read: multiple thousand) on NSQ as the original EMDR message contains many rowsets. Identical rowsets get deduplicated. This leads to way less messages on the queue as most of the time most of the market does not change between bulk updates from EMDR. In order for this feature to work properly it is best to only consume messages produced by a specific generator (you can set filters). Most generators have multiple instances running all over the world, so filtering for a specific one (say CRESTMarketTrawler v0.5.0) should not greatly reduce the reliability. Of course you can also choose to set no flters at all.

## Installation
Either use the prebuilt Docker images and pass the appropriate env vars (see below), or:

* Clone this repo into your gopath
* Run `go get`
* Run `go build`


## Deployment Info
Builds and releases are handled by Drone.

Environment Variable | Default | Description
--- | --- | ---
LOG_LEVEL | info | Threshold for logging messages to be printed
EMDR_RELAY_URL | tcp://relay-eu-germany-1.eve-emdr.com:8050 | EMDR relay to connect to
NSQ_URL | nsqd:4150 | Hostname/IP of the NSQD instance to connect to
GENERATOR_NAME | none / match all | Only forward messages by a generator whose name matches this regex (see input example below). Remember to properly escape special characters in the regex.
GENERATOR_VERSION | none / match all | Only forward messages by a generator whose version matches this regex (see input example below). Remember to properly escape special characters in the regex.
CACHE_PATH | cache.db | Path to persistent deduplication cache

## Todo
- [ ] General code cleanup (this is my first Go project), add performance metrics
- [ ] Tests would be nice

## Example Message

Input from EMDR:
```json
{
  "resultType" : "orders",
  "version" : "0.1",
  "uploadKeys" : [
    { "name" : "emk", "key" : "abc" },
    { "name" : "ec" , "key" : "def" }
  ],
  "generator" : { "name" : "Yapeal", "version" : "11.335.1737" },
  "currentTime" : "2011-10-22T15:46:00+00:00",
  "columns" : ["price","volRemaining","range","orderID","volEntered","minVolume","bid","issueDate","duration","stationID","solarSystemID"],
  "rowsets" : [
    {
      "generatedAt" : "2011-10-22T15:43:00+00:00",
      "regionID" : 10000065,
      "typeID" : 11134,
      "rows" : [
        [8999,1,32767,2363806077,1,1,false,"2011-12-03T08:10:59+00:00",90,60008692,30005038],
        [11499.99,10,32767,2363915657,10,1,false,"2011-12-03T10:53:26+00:00",90,60006970,null],
        [11500,48,32767,2363413004,50,1,false,"2011-12-02T22:44:01+00:00",90,60006967,30005039]
      ]
    },
    {
      "generatedAt" : "2011-10-22T15:42:00+00:00",
      "regionID" : null,
      "typeID" : 11135,
      "rows" : [
        [8999,1,32767,2363806077,1,1,false,"2011-12-03T08:10:59+00:00",90,60008692,30005038],
        [11499.99,10,32767,2363915657,10,1,false,"2011-12-03T10:53:26+00:00",90,60006970,null],
        [11500,48,32767,2363413004,50,1,false,"2011-12-02T22:44:01+00:00",90,60006967,30005039]
      ]
    },
    {
      "generatedAt" : "2011-10-22T15:43:00+00:00",
      "regionID" : 10000067,
      "typeID" : 11136,
      "rows" : []
    }
  ]
}
```

Output to NSQ (multiple messages):
```json
[
  {
    "typeID": 11134,
    "regionID": 10000065,
    "orders": [
      {
        "volRemaining": 1,
        "volEntered": 1,
        "stationID": 60008692,
        "solarSystemID": 30005038,
        "range": 32767,
        "price": 8999,
        "orderID": 2363806077,
        "minVolume": 1,
        "issueDate": "2011-12-03T08:10:59+00:00",
        "duration": 90,
        "bid": false
      },
      {
        "volRemaining": 10,
        "volEntered": 10,
        "stationID": 60006970,
        "solarSystemID": null,
        "range": 32767,
        "price": 11499.99,
        "orderID": 2363915657,
        "minVolume": 1,
        "issueDate": "2011-12-03T10:53:26+00:00",
        "duration": 90,
        "bid": false
      },
      {
        "volRemaining": 48,
        "volEntered": 50,
        "stationID": 60006967,
        "solarSystemID": 30005039,
        "range": 32767,
        "price": 11500,
        "orderID": 2363413004,
        "minVolume": 1,
        "issueDate": "2011-12-02T22:44:01+00:00",
        "duration": 90,
        "bid": false
      }
    ],
    "generatedAt": "2011-10-22T15:43:00+00:00"
  },
  {
    "typeID": 11135,
    "regionID": null,
    "orders": [
      {
        "volRemaining": 1,
        "volEntered": 1,
        "stationID": 60008692,
        "solarSystemID": 30005038,
        "range": 32767,
        "price": 8999,
        "orderID": 2363806077,
        "minVolume": 1,
        "issueDate": "2011-12-03T08:10:59+00:00",
        "duration": 90,
        "bid": false
      },
      {
        "volRemaining": 10,
        "volEntered": 10,
        "stationID": 60006970,
        "solarSystemID": null,
        "range": 32767,
        "price": 11499.99,
        "orderID": 2363915657,
        "minVolume": 1,
        "issueDate": "2011-12-03T10:53:26+00:00",
        "duration": 90,
        "bid": false
      },
      {
        "volRemaining": 48,
        "volEntered": 50,
        "stationID": 60006967,
        "solarSystemID": 30005039,
        "range": 32767,
        "price": 11500,
        "orderID": 2363413004,
        "minVolume": 1,
        "issueDate": "2011-12-02T22:44:01+00:00",
        "duration": 90,
        "bid": false
      }
    ],
    "generatedAt": "2011-10-22T15:42:00+00:00"
  },
  {
    "typeID": 11136,
    "regionID": 10000067,
    "orders": [

    ],
    "generatedAt": "2011-10-22T15:43:00+00:00"
  }
]
```
