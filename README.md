# go-logger-syslog
logging our golang application using rsyslog + logstash + kibana 

## Logging in our golang application
#### Creating our logger
```
log, err := logger.NewLoggerFromDSN(loggerDSN, appName, *env)
```

The DSN (defined in our `config.toml`) should looks like:
```
dsn = "kibana://?level=debug"
```

Schema could be:
1. kibana `kibana://?level=debug`
1. stdout `stdout://?level=info`
1. discardall `discardall://?level=debug`

Level parameter could be `debug` or `info`

#### Using our logger
```
log.Info("Running...", logger.NewField("newField1", "value1"), logger.NewField("newField2", 2))
```

```
log.Debug("Debugging Running...", logger.NewField("newField1", "value1"), logger.NewField("newField2", 2))
```


## Instaling rsyslog and logstash in (for Red Hat-based systems)
#### Configure Rsyslog

Active mmjsonparse it will be used to parse all the `@cee json messages
```
module(load="mmjsonparse")
```

Create rulset that uses `mmjsonparse` to parse the `@cee messages and do action to forward to logstash udp port the json message
```
ruleset(name="remoteAllJsonLog") {
    action(type="mmjsonparse")
    if $parsesuccess == "OK" then {
        action(
            type="omfwd"
            Target="localhost"
            Port="5514"
            Protocol="udp"
            template="allJsonLogTemplate"
        )
    }
    stop
}
```

Define a template to get all json fields of the message
```
template(name="allJsonLogTemplate" type="list") {
    property(name="$!all-json")
}
```

Get all the logs using this udp port and use the ruleset `remoteAllJsonLog`
```
module(load="imudp") # needs to be done just once
input(type="imudp" port="514" ruleset="remoteAllJsonLog")
```

#### Configure logstash
Input get messages from port 5514 UDP
```
input {
  udp {
      port => 5514
      codec => json
  }
}
```

Output send messages to elasticsearch
```
output {
    elasticsearch {
        hosts => "[ELASTICSEARCH_HOST]"
    }
}
```