syntax = "proto3";
package pb;

service CommunicateService {
    rpc MessageChannel (stream Message) returns (stream Message) {}
}
message Message {
    oneof Message {
        Example1 example1= 1;
        Example2 example1= 2;
    }
}
message Example1 {
    string ttype = 1;
}
message Example2 {
    string ttype = 1;
}