# AuditLogParser
[build](https://github.com/kenjords/auditlogparser/actions/workflows/test_and_build.yml/badge.svg) [release](https://github.com/kenjords/auditlogparser/actions/workflows/release_build.yml/badge.svg)

AuditLogParser is a tool for parsing DSE audit.log files and formatting them into consumable JSON.

## Installation
Download the latest release from the [releases page](https://github.com/kenjords/AuditLogParser/releases) and extract  
the contents to a directory of your choice.  
```bash
tar -C /path/to/bin -xzf auditlogparser_<version>_<os>_<arch>.tar.gz
```

## Usage
Usage is pretty simple. 

```bash
auditlogparser -file /path/to/audit.log
```
This will output the audit log in JSON format to stdout. 
You can then use tools such as `jq` to filter the output. 

```bash
auditlogparser -file /path/to/audit.log | jq '. | select(.operation == "SELECT")'
```

## Options

| Option | Description |
| ------ | ----------- |
| -file | Path to the audit.log file to parse. |

