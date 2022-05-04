# grpc_demo

```
protoc -I pb/ pb/cyber.proto --go_out=plugins=grpc:service
python -m grpc_tools.protoc -I pb/ --python_out=example/python/ --grpc_python_out=example/python/ pb/cyber.proto
