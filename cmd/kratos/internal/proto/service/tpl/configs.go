package tpl

var ConfigYaml = `trace:
  endpoint: http://127.0.0.1:14268/api/traces
server:
  server_id: 1
  http:
    addr: 0.0.0.0:0
    timeout: 10s
  grpc:
    addr: 0.0.0.0:0
    timeout: 10s
data:
  database:
    driver: mysql
    source: root:123456@tcp(192.168.8.91:3306)/blind_box?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&multiStatements=true
    log_level: 3
  redis:
    addr: 192.168.8.91:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
    database: 0
  kafka:
    addrs:
      - 127.0.0.1:9092
`

var RegistryYaml = `consul:
  address: 192.168.8.91:8500
  scheme: http`
