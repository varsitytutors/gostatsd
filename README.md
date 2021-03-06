gostatsd
========

[![Godoc](https://godoc.org/github.com/atlassian/gostatsd?status.svg)](https://godoc.org/github.com/atlassian/gostatsd)
[![Build Status](https://app.codeship.com/projects/208ba270-3165-0137-d5c7-3ed4c607a748/status?branch=master)]
[![license](https://img.shields.io/github/license/atlassian/gostatsd.svg)](https://github.com/atlassian/gostatsd/blob/master/LICENSE)

An implementation of [Etsy's][etsy] [statsd][statsd] in Go,
based on original code from [@kisielk](https://github.com/kisielk/).

The project provides both a server called "gostatsd" which works much like
Etsy's version, but also provides a library for developing customized servers.

Backends are pluggable and only need to support the [backend interface](backend.go).

Being written in Go, it is able to use all cores which makes it easy to scale up the
server based on load. The server can also be run HA and be scaled out, see
[Load balancing and scaling out](https://github.com/atlassian/gostatsd#load-balancing-and-scaling-out).

Building the server
-------------------
Gostatsd currently targets Go 1.10.2.  There are no known hard dependencies in the code beween 1.9 and 1.10.2, but some may be introduced in future.

From the `gostatsd` directory run `make build`. The binary will be built in `build/bin/<arch>/gostatsd`.

You will need to install the Golang build dependencies by running `make setup` in the `gostatsd` directory. This must be done before the first build,
and again if the dependencies change.  A [protobuf](https://github.com/protocolbuffers/protobuf) installation is expected to be found in the `tools/`
directory.  Managing this in a platform agnostic way is difficult, but PRs are welcome. Hopefully it will be sufficient to use the generated protobuf
files in the majority of cases.

If you are unable to build `gostatsd` please try running `make setup` again before reporting a bug.

Running the server
------------------
`gostatsd --help` gives a complete description of available options and their
defaults. You can use `make run` to run the server with just the `stdout` backend
to display info on screen.

You can also run through `docker` by running `make run-docker` which will use `docker-compose`
to run `gostatsd` with a graphite backend and a grafana dashboard.

While not generally tested on Windows, it should work.  Maximum throughput is likely to be better on
a linux system, however.

Configuring the server mode
---------------------------
The server can currently run in two modes: `standalone` and `forwarder`.  It is configured through the top level
`server-mode` configuration setting.  The default is `standalone`.

In `standalone` mode, raw metrics are processed and aggregated as normal, and aggregated data is submitted to
configured backends (see below)

In `forwarder` mode, raw metrics are collected from a frontend, and instead of being aggregated they are sent via http
to another gostatsd server after passing through the processing pipeline (cloud provider, static tags, filtering, etc).

A `forwarder` server is intended to run on-host and collect metrics, forwarding them on to a central aggregation
service.  At present the central aggregation service can only scale vertically, but horizontal scaling through
clustering is planned.

Configuring `forwarder` mode requires a configuration file, with a section named `http-transport`.  The raw version
spoken is not configurable per server (see [HTTP.md] for version guarantees).  The configuration section allows the
following configuration options:

- `client-timeout`: duration for the http client timeout.  Defaults to `10s`
- `compress`: boolean indicating if the payload should be compressed.  Defaults to `true`
- `enable-http2`: boolean to enable the usage of http2 on the request.  There seems to be some incompatibility with the
  golang http2 implementation and AWS ELB/ALBs.  If you experience strange timeouts and hangs, this should be the first
  thing to disable.  Defaults to `false`
- `api-endpoint`: configures the endpoint to submit raw metrics to.  This setting should be just a base URL, for example
  `https://statsd-aggregator.private`, with no path.  Required, no default
- `max-requests`: maximum number of requests in flight.  Defaults to `1000` (which is probably too high)
- `max-request-elapsed-time`: duration for the maximum amount of time to try submitting data before giving up.  This
  includes retries.  Defaults to `30s` (which is probably too high)
- `network`: the network type to use, probably `tcp`, `tcp4`, or `tcp6`.  Defaults to `tcp`
- `consolidator-slots`: number of slots in the metric consolidator.  Memory usage is a function of this.  Lower values
  may cause blocking in the pipeline (back pressure).  A UDP only receiver will never use more than the number of
  configured parsers (`--max-parsers` option).  Defaults to the value of `--max-parsers`, but may require tuning for
  HTTP based servers.
- `flush-interval`: duration for how long to batch metrics before flushing. Should be an order of magnitude less than
  the upstream flush interval. Defaults to `1s`

Configuring HTTP servers
------------------------
The service supports multiple HTTP servers, with different configurations for different requirements.  All http servers
are named in the top level `http-servers` setting.  It should be a space separated list of names.  Each server is then
configured by creating a section in the configuration file named `http.<servername>`.  An http server section has the
following configuration options:

- `address`: the address to bind to
- `enable-prof`: boolean indicating if profiler endpoints should be enabled. Default `false`
- `enable-expvar`: boolean indicating if expvar endpoints should be enabled. Default `false`
- `enable-ingestion`: boolean indicating if ingestion should be enabled. Default `false`
- `enable-healthcheck`: boolean indicating if healthchecks should be enabled. Default `true`

For example, to configure a server with a localhost only diagnostics endpoint, and a regular ingestion endpoint that
can sit behind an ELB, the following configuration could be used:

```config.toml
backends='stdout'
http-servers='receiver profiler'

[http.receiver]
address='0.0.0.0:8080'
enable-ingestion=true

[http.profiler]
address='127.0.0.1:6060'
enable-expvar=true
enable-prof=true
```

There is no capability to run an https server at this point in time, and no auth (which is why you might want different
addresses).  You could also put a reverse proxy in front of the service.  Documentation for the endpoints can be found
under HTTP.md

Configuring backends and cloud providers
----------------------------------------
Backends and cloud providers are configured using `toml`, `json` or `yaml` configuration file
passed via the `--config-path` flag. For all configuration options see source code of the backends you
are interested in. A cloudprovider should not be used on the aggregation server when forwarding data to
it, as the source IP address is not propagated.  A cloudprovider can be used on the forwarder host, however.
Configuration file might look like this:
```
[graphite]
	address = "192.168.99.100:2003"

[datadog]
	api_key = "my-secret-key" # Datadog API key required.

[statsdaemon]
	address = "docker.local:8125"
	disable_tags = false

[aws]
	max_retries = 4

[newrelic]
	address = "http://localhost:8001/v1/data"
	event-type = "GoStatsD"
	#see full configuration options further below
```

New Relic Backend
-----------------------------
Supports two routes for flushing metrics to New Relic.
- Directly to the Insights Collector - [Insights Event API](https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/send-custom-events-event-api)
- Via the Infrastructure Agent's inbuilt HTTP Server

### [New Relic Insights Event API](https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/send-custom-events-event-api)
Sending directly to the Event API alleviates the requirement of needing to have the New Relic Infrastructure Agent. Therefore you can run this from nearly anywhere for maximum flexibility. This also becomes a shorter data path with less resource requirements becoming a simpler setup.

To use this method, create an Insert API Key from here: https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/manage/api_keys

```
#Example configuration

[newrelic]
    address = "https://insights-collector.newrelic.com/v1/accounts/YOUR_ACCOUNT_ID/events"
    api-key = "yourEventAPIInsertKey"
```

### [New Relic Infrastructure Agent](https://newrelic.com/products/infrastructure)
Sending via the Infrastructure Agent's inbuilt HTTP server provides additional features, such as automatically applying additional metadata to the event the host may have such as AWS tags, instance type, host information, labels etc.

The payload structure required to be accepted by the agent can be viewed [here.](https://github.com/newrelic/infra-integrations-sdk/blob/master/docs/v2tov3.md#v2-json-full-sample)

To enable the HTTP server, modify /etc/newrelic.yml to include the below, and restart the agent ([Step 1.2](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/statsd-monitoring-integration#install)).
```
http_server_enabled: true
http_server_host: 127.0.0.1 #(default host)
http_server_port: 8001 #(default port)
```

Additional options are available to rename attributes if required.
```
[newrelic]
	tag-prefix = ""
	metric-name = "name"
	metric-type = "type"
	per-second = "per_second"
	value = "value"
	timer-min = "min"
	timer-max = "max"
	timer-count = "samples_count"
	timer-mean = "samples_mean"
	timer-median = "samples_median"
	timer-stddev = "samples_std_dev"
	timer-sum = "samples_sum"
	timer-sumsquare = "samples_sum_squares"
```


Configuring timer sub-metrics
-----------------------------
By default, timer metrics will result in aggregated metrics of the form (exact name varies by backend):
```
<base>.Count
<base>.CountPerSecond
<base>.Mean
<base>.Median
<base>.Lower
<base>.Upper
<base>.StdDev
<base>.Sum
<base>.SumSquares
```


In addition, the following aggregated metrics will be emitted for each configured percentile:
```
<base>.Count_XX
<base>.Mean_XX
<base>.Sum_XX
<base>.SumSquares_XX
<base>.Upper_XX - for positive only
<base>.Lower_-XX - for negative only
```


These can be controlled through the `disabled-sub-metrics` configuration section:
```
[disabled-sub-metrics]
# Regular metrics
count=false
count-per-second=false
mean=false
median=false
lower=false
upper=false
stddev=false
sum=false
sum-squares=false

# Percentile metrics
count-pct=false
mean-pct=false
sum-pct=false
sum-squares-pct=false
lower-pct=false
upper-pct=false
```


By default (for compatibility), they are all false and the metrics will be emitted.



Sending metrics
---------------
The server listens for UDP packets on the address given by the `--metrics-addr` flag,
aggregates them, then sends them to the backend servers given by the `--backends`
flag (space separated list of backend names).

Currently supported backends are:

* graphite
* datadog
* statsdaemon
* stdout
* cloudwatch
* newrelic

The format of each metric is:

    <bucket name>:<value>|<type>\n

* `<bucket name>` is a string like `abc.def.g`, just like a graphite bucket name
* `<value>` is a string representation of a floating point number
* `<type>` is one of `c`, `g`, or `ms` for "counter", "gauge", and "timer"
respectively.

A single packet can contain multiple metrics, each ending with a newline.

Optionally, `gostatsd` supports sample rates (for simple counters, and for timer counters) and tags:

* `<bucket name>:<value>|c|@<sample rate>\n` where `sample rate` is a float between 0 and 1
* `<bucket name>:<value>|c|@<sample rate>|#<tags>\n` where `tags` is a comma separated list of tags
* `<bucket name>:<value>|<type>|#<tags>\n` where `tags` is a comma separated list of tags

Tags format is: `simple` or `key:value`.


A simple way to test your installation or send metrics from a script is to use
`echo` and the [netcat][netcat] utility `nc`:

    echo 'abc.def.g:10|c' | nc -w1 -u localhost 8125

Monitoring
----------
Many metrics for the internal processes are emitted.  See METRICS.md for details.  Go expvar is also
exposed if the `--profile` flag is used.

Memory allocation for read buffers
----------------------------------
By default `gostatsd` will batch read multiple packets to optimise read performance. The amount of memory allocated
for these read buffers is determined by the config options:

    max-readers * receive-batch-size * 64KB (max packet size)

The metric `avg_packets_in_batch` can be used to track the average number of datagrams received per batch, and the
`--receive-batch-size` flag used to tune it.  There may be some benefit to tuning the `--max-readers` flag as well.

Using the library
-----------------
In your source code:

    import "github.com/atlassian/gostatsd/pkg/statsd"

Documentation can be found via `go doc github.com/atlassian/gostatsd/pkg/statsd` or at
https://godoc.org/github.com/atlassian/gostatsd/pkg/statsd

Contributors
------------

Pull requests, issues and comments welcome. For pull requests:

* Add tests for new features and bug fixes
* Follow the existing style
* Separate unrelated changes into multiple pull requests

See the existing issues for things to start contributing.

For bigger changes, make sure you start a discussion first by creating an issue and explaining the intended change.

Atlassian requires contributors to sign a Contributor License Agreement, known as a CLA. This serves as a record stating that the contributor is entitled to contribute the code/documentation/translation to the project and is willing to have it used in distributions and derivative works (or is willing to transfer ownership).

Prior to accepting your contributions we ask that you please follow the appropriate link below to digitally sign the CLA. The Corporate CLA is for those who are contributing as a member of an organization and the individual CLA is for those contributing as an individual.

* [CLA for corporate contributors](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=e1c17c66-ca4d-4aab-a953-2c231af4a20b)
* [CLA for individuals](https://na2.docusign.net/Member/PowerFormSigning.aspx?PowerFormId=3f94fbdc-2fbe-46ac-b14c-5d152700ae5d)

License
-------

Copyright (c) 2012 Kamil Kisiel.
Copyright @ 2016-2017 Atlassian Pty Ltd and others.

Licensed under the MIT license. See LICENSE file.

[etsy]: https://www.etsy.com
[statsd]: https://www.github.com/etsy/statsd
[netcat]: http://netcat.sourceforge.net/
