package tpl

var ConfTemplate = `syntax = "proto3";

package {{ .ServiceLower }}.{{ .Mode }}.conf;

option go_package = "gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf;conf";

import "google/protobuf/duration.proto";

message Bootstrap {
  Trace trace = 1;
  Server server = 2;
  Service service = 3;
  Data data = 4;
}

message Trace {
  string endpoint = 1;
}

message Server {
  message HTTP {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration timeout = 3;
  }
  message GRPC {
    string network = 1;
    string addr = 2;
    google.protobuf.Duration timeout = 3;
  }
  HTTP http = 1;
  GRPC grpc = 2;
}

message Service {
  int64 server_id = 1;

  message Aes {
    string key = 1;
    string vi = 2;
  }
  Aes aes = 2;
  message Job {
    int32 enable = 1;
    google.protobuf.Duration cycle_time = 2;
  }
  Job job = 3;
}


message Data {
  message Database {
    string driver = 1;
    string source = 2;
    int64 log_level = 3;
  }
  message Redis {
    string network = 1;
    string addr = 2;
    string password = 3;
    int32 database = 4;
    int32 pool_size = 5;
    google.protobuf.Duration read_timeout = 7;
    google.protobuf.Duration write_timeout = 8;
  }
  message Kafka {
    repeated string addrs = 1;
  }
  Database database = 1;
  Redis redis = 2;
  Kafka kafka = 3;
}

message Registry {
  message Consul {
    string address = 1;
    string scheme = 2;
  }
  Consul consul = 1;
}
`
