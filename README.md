<p align="center">
  <img alt="Logo" src="https://user-images.githubusercontent.com/48660576/91569625-85384600-e94e-11ea-8ce2-5ede821d5d94.png" height="200" />
  <h3 align="center"><strong>A terminal UI to monitor and query Elasticsearch. </strong></h3>
</p>
<p align="center">
   <a href="https://goreportcard.com/badge/github.com/serdarkkts/qelk"><img src="https://goreportcard.com/badge/github.com/serdarkkts/qelk?" alt="Go Report"></a>
  <a href="https://github.com/serdarkkts/qelk/blob/master/LICENSE"><img src="https://img.shields.io/github/license/serdarkkts/qelk?" alt="Apache 2.0"></a>
  <a href="https://asciinema.org/a/WGDybzDPBY7yzemNvmPRUSiu2"><img src="https://img.shields.io/badge/demo-terminal%20session-brightgreen" alt="Demo"></a>
<a href="https://github.com/serdarkkts/qelk/releases"><img src="https://img.shields.io/github/v/release/serdarkkts/qelk" alt="Demo"></a>
  
</p>

___

## What is qelk?

Qelk helps you to monitor your Elasticsearch cluster, execute search queries
and creates dashboards for your search queries from your command line. 

Key features are: 

* Colorful Terminal UI.
* Node monitoring( Memory usage, CPU, JVM, Disk usage...) 
* Node Filtering
* Custom dashboards
* Query parameter search 
* Search results with custom format
* Supports only ES 7.x

## Why qelk?

Get things done with a few commands from the cli.


## Installation

* Go
```bash
$ go get -u github.com/serdarkkts/qelk
```
* Go to the [releases](https://github.com/serdarkkts/qelk/releases) page and download the latest one for your platform.


## Configuration

Before using the qelk, create a configuration file `.qelk.yaml` in your home directory. After creating the file, there are three main configurations you can make;

### Urls
Provide your ES instance urls. You can provide more than one. If you leave this field empty, qelk will attempt to connect ES instance at localhost:9200.

An example configuration:
```yaml
urls:
  - http://localhost:9200
  - http://localhost:9300
  - http://localhost:9220
```

### Authentication
Provide your ES instance credentials. You can provide both username and password. If you don't want to keep your password on the config file for security reasons, you can provide only username, password will be prompted.

An example configuration:
```yaml
username: elastic
password: elastic  #leave empty for prompt
```

### Dashboards

Define your custom dashboards. To create a custom dashboard, firstly add a dashboards key in your config file like this.

```yaml
dashboards:
```

After that, add a key with your dashboard name under dashboards key like this.

```yaml
dashboards:
  nginx:
```

After that, there are 5 fields you can define in your custom dashboard. Providing all fields is not necessary, you can define only the fields you need.

* index (string) - Define indices for this dashboard. Regex supported. 
* format (array) - Define format for results. 
* sort (array)   - Define a sorting for this dashboard.
* size (int) - Define size of this dashboard. 
* query (string) - Define query for this dashboard.

An example configuration with 2 custom dashboards:
```yaml
dashboards:
  nginx:
    index: logstash-*
    format:
      - "@timestamp"
      - logLevel
      - message
    sort:
      - "@timestamp:desc"
    size: 50
    query: logLevel:ERROR
  products:
    index: products*
    format:
      - created
      - name
      - description
    sort:
      - created:desc
    size: 12
    query: Onions
```

## Usage

The application has 3 main commands. These are `search`, `inspect`, and `stats`. You can run `stats` and `inspect` commands from the UI, but the `search` command can only be run from the command line.

### Search Command

Search command performs a search request.

```bash
Usage:
  qelk search [flags]

Flags:
  -f, --format strings   specify the format
  -i, --index string     specify the index (default "*")
  -q, --query string     specify the query
  -n, --size int         specify the size (default 15)
  -s, --sort strings     specify the sorting filter

```

Example:

```bash
$ qelk search -q name:Onions -n 10 -f created,name,description -s created:desc
```

### Inspect Command

Inspect your custom dashboards. Inspect command requires your custom dashboard name as argument.

For example, if you defined your dashboard in config file like this:

```yaml
dashboards:
  nginx:
    index: logstash-*
    format:
      - "@timestamp"
      - logLevel
      - message
    sort:
      - "@timestamp:desc"
    size: 50
    query: logLevel:ERROR
```

you can inspect your dashboard like this:

```bash
$ qelk inspect nginx
```

### Stats Command

Monitor your ES cluster. Fields in the Stats interface mean the following.



| Name             | Description                               |
|------------------|-----------------------------------------------------------------|
| `CLUSTER`   | Name of the cluster, based on the `cluster.name` setting.                       |
| `STATUS`           | Health status of the cluster. |
| `CPU`        |  Percentage of CPU used across all selected nodes. `Returns -1` if not supported. |
| `CPU.Avaible Proc.`   | Number of processors available to JVM across all selected nodes. |
| `CPU.Allocated Proc.` | Number of processors used to calculate thread pool size across all selected nodes.|
| `MEMORY`        | Percentage of physical memory in use across all selected nodes.|
| `MEMORY.Free`   | Amount of free physical memory across all selected nodes.       |
| `MEMORY.Used`           | Amount of physical memory in use across all selected nodes.  |
| `MEMORY.Total`        | Total amount of physical memory across all selected nodes. |
| `DISK.Avaible`           | Total amount of disk space available to JVM in file stores across all selected nodes. |
| `DISK.Free`        | Amount of unallocated disk space in file stores across all selected nodes. |
| `DISK.Total`           | Total size of all file stores across all selected nodes.  |
| `NODES.Total`           | Total number of nodes selected by the request.           |
| `NODES.Successful`        | Number of nodes that responded successfully to the request. |
| `NODES.Failed`           |  Number of nodes that rejected the request or failed to respond.   |
| `INDICES.Total`           | Total number of indices with shards assigned to selected nodes.      |
| `INDICES.Documents`        | Total number of non-deleted documents across all primary shards assigned to selected nodes. |
| `INDICES.Uptime`           | Uptime duration since JVM last started.          |
| `JVM.Used`           | Memory currently in use by the heap across all selected nodes.     |
| `JVM.Max`        | Maximum amount of memory, in bytes, available for use by the heap across all selected nodes. |
| `JVM.Threads`           | Number of active threads in use by JVM across all selected nodes.   |
| `SHARDS.Total`        |  Total number of shards assigned to selected nodes. |
| `SHARDS.Primary`           | Number of primary shards assigned to selected nodes.          |
| `SHARDS.Replication`        | Ratio of replica shards to primary shards across all selected nodes. |
| `SHARDS.Size`        | Total size of all shards assigned to selected nodes. |
| `VERSION`           |  Contains statistics about the JVM versions used by selected nodes. |
| `FILTER`        | Selected node filter.  |

You can run stats command like this,

```bash
$ qelk stats
```

## UI
Stats and inspect commands can be run through the UI as follows:

First, run the following command.
```bash
$ qelk 
```

After that you can press <kbd>F1</kbd> for `stats`, <kbd>F2</kbd> for `inspect`. You can use <kbd>tab</kbd> key for moving between different widgets, <kbd>ENTER</kbd> key for pressing buttons.

## TODO
Want to contribute? Here is a short TODO list of things that aren't implemented in this application that would be nice:

* [CACert connection](https://www.elastic.co/guide/en/elasticsearch/reference/master/encrypting-communications-certificates.html).

* Add support for other ES versions.

* Add themes for UI.

* Add [shell completions](https://github.com/spf13/cobra/blob/master/shell_completions.md). 


## License

This project is under the [Apache License 2.0](https://github.com/serdarkkts/qelk/blob/master/LICENSE)
