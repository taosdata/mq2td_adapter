# Receive mqtt message and write TDengine
Access the MQTT message as a consumer, parse the json according to the configuration and write it into TDengine

English | [简体中文](README-CN.md)

## Compilation method
1. Install golang(1.14+) [https://golang.google.cn/doc/install](https://golang.google.cn/doc/install)
2. Set up golang proxy
```shell
go env -w GOPROXY=https://goproxy.cn,direct
```
3. Install the TDengine client [https://www.taosdata.com/en/getting-started/](https://www.taosdata.com/en/getting-started/)
4. Configure C Compiler Environment(GCC)
5. Execute in the project directory
```shell
go build
```

## MQTT configuration

```json
{
  "address": "mqtt Address",
  "clientID": "If the client ID is not set, use uuid",
  "username": "Username",
  "password": "Password",
  "keepAlive": 30,
  "caPath": "ca certification path",
  "certPath": "Certificate Path",
  "keyPath": "Certificate key path"
}
```
`keepAlive`: Alive time (in seconds)

## TDengine configuration

```json
{
  "host": "Address",
  "port": 6030,
  "user": "Username",
  "password": "Password",
  "db": "Database"
}
```

`port` is the TDengine service port

## Parsing rule configuration

```json
[
  {
    "rule_name": "Rule Name",
    "topic": "topic",
    "rule": {
      "s_table": "Corresponding to STable name",
      "table": {
        "default_value": "Default value",
        "path": "json path"
      },
      "tags": [
        {
          "name": "corresponding to the tag name in TDengine",
          "value_type": "Type of value",
          "length": "Maximum length of the value (need to be set when the value type is string)",
          "default_value": "Default value",
          "path": "json path",
          "time_layout": "Time formatted layout (need to be set when the value type is timeString)"
        }
      ],
      "columns": [
        {
          "name": "corresponding to the column name in TDengine",
          "value_type": "Type of value",
          "length": "Maximum length of the value (need to be set when the value type is string)",
          "default_value": "Default value",
          "path": "json path",
          "time_layout": "Time formatted layout (need to be set when the value type is timeString)"
        }
      ]
    }
  }
]
```
* Default value: Use the default value when the value corresponding to path is not found in json
* Type of value: `"int"
  "float"
  "bool"
  "string"
  "timeString"
  "timeSecond"
  "timeMillisecond"
  "timeMicrosecond"
  "timeNanosecond"`
* json path see [https://github.com/tidwall/gjson](https://github.com/tidwall/gjson)
* When the value type is `string`, the `length` parameter must be set
* `tags` set at least one
* `columns` set at least two, the first parameter name must be `ts` and the type must be one of `"timeString" "timeSecond" "timeMillisecond" "timeMicrosecond" "timeNanosecond"`
  `timeLayout` is golang time formatting template [https://golang.google.cn/pkg/time/#pkg-constants](https://golang.google.cn/pkg/time/#pkg-constants)
## Parameters
```
  -c string
        Configuration file path (default "./config/config.json")
  --rc string
        Rule configuration file path (default "./config/rule.json")
```

## Log configuration

```json
{
   "log": {
     "path": "/var/log/taos2",
     "rotationCount": 7,
     "rotationTime": "1d",
     "rotationSize": "1GB"
   }
}
```

* `log.path` log file directory
* `log.rotationCount` the number of logs to keep
* `log.rotationTime` log splitting time limit
* `log.rotationSize` log splitting size limit

When either the time limit or size limit is triggered, log splitting will be performed. If the number of log files is greater than the number of retained count, the oldest log file will be deleted according to the modification time.

## Configuration example
See [example folder](example)

## Launch

Assume you are under the directory where `mq2td_adapter` is located, please use the command below to run:

```shell
nohup ./mq2td_adapter >/dev/null 2>&1 &
```
