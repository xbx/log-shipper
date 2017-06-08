# log-shipper ALPHA

Nginx to Elastichsearch Log Shipper

Log Shipper is a simple tool to `tail -f` NGiNX logs, parse them to `json` format and insert them into a Elasticsearch.
This is an Alpha version which just works but need some improvements.

![](https://github.com/xbx/log-shipper/raw/master/log-shipper.png)

## Usage

    log-shipper -elastichost some-elasticsearch /var/log/nginx/*log

## Todo

* Configurable Elasticsearch port
