# Json RPC interceptor

The tool allows to intercept HTTP (Json RPC) requests and publish the results as Prometheus metrics.

## How it works
1. The interceptor is Nginx proxy server that listens on a specific port and forwards the requests to the JRPC app.
2. Nginx contains Lua script that parses the requests to find out JRPC methods.
3. Nginx sends the logs and other metrics (e.g. response time and parsed JRPC method) to syslog server in specific format.
4. Syslog server is Go application that parses that logs and publishes the metrics to Prometheus.

Your JRPC app and the interceptor should be running on the same machine (network).


## How to run
1. Build the interceptor docker image:
```bash
docker build -t jrpc-interceptor .
```
2. Run your JRPC app, e.g.:
```bash
docker run --net=host -it nethermindeth/nethermind:latest
```
3. Run the interceptor:
```bash
docker run --net="host" -e LISTEN_PORT=0.0.0.0:8081 -e SERVICE_TO_PROXY=0.0.0.0:8545 -e LOG_SERVER_URL=0.0.0.0:514 jrpc-interceptor
```

Where:
- `LISTEN_PORT` - the port / IP where the interceptor listens for incoming requests.
- `SERVICE_TO_PROXY` - the IP / port of the JRPC app.
- `LOG_SERVER_URL` - the IP / port of the syslog server.
- `PROMETHEUS_URL` - the IP / port of the Prometheus server. Optional, "0.0.0.0:9100" by default.
- `USE_PROMETHEUS` - whether to publish metrics to Prometheus. Optional, "true" by default.
- `LOG_SERVER_DEBUG` - whether to print the logs to stdout. Optional, "true" by default.

4. Send requests to the interceptor:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://0.0.0.0:8081
```
5. Check the logs. You should see something like this:
```bash
2024/08/13 09:01:37 ef01982f3841: 192.168.65.1|http|localhost|POST|HTTP/1.1|/|200|0.002|193|326|0.003|eth_blockNumber
```
6. Open URL `http://localhost:9100/metrics` in your browser to see the metrics.

## Local development
1. Add download dependencies:
```bash
go mod download
```
2. Build the package
```bash
go build .
```
3. Run the interceptor:
```bash
./jrpc-interceptor -debug=${LOG_SERVER_DEBUG:-true} -listenSyslog=${LOG_SERVER_URL:-"0.0.0.0:514"} -listenHTTP=${PROMETHEUS_URL:-"0.0.0.0:9100"} -usePrometheus=${USE_PROMETHEUS:-true}
```

### Download docker image from github container registry

1. Create a [personal access token ](https://github.com/settings/tokens)(PAT) on github with repo access
2. Set the PAT as an environment variable

```
export CR_PAT=<your_pat>
```
3. Login to github container registry

```
echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
```

4. Download the image

```
docker pull ghcr.io/NethermindEth/jrpc-interceptor:main
```

## Nginx configuration
All the configuration can be found inside the `nginx` folder.

Basic info:

- `worker_processes` sets to 4;
- `worker_connections` sets 2048;
- `client_body_buffer_size` sets to `10M`;;
- `client_max_body_size` sets to `1000M` due to the large size of some jrpc requests;

Used `openresty/openresty:alpine` to be able to use Lua scripts.

### Syslog format
The interceptor sends the logs to the syslog server in the following format:
```
$remote_addr|$scheme|$host|$request_method|$server_protocol|$request_uri|$status|$request_time|$request_length|$bytes_sent|$upstream_response_time|$jrpc_method
```
That returns something like
```
2024/08/13 09:01:37 ef01982f3841: 192.168.65.1|http|localhost|POST|HTTP/1.1|/|200|0.002|193|326|0.003|eth_blockNumber
```

### Limitations and workarounds
- All the original headers passes as is, except `Host` header. It's replaced with the `SERVICE_TO_PROXY` value.

## Errors and Troubleshooting
Instead of a valid `jrpc_method` you can show following errors:
- `no_method_field` this means that original request contains `jsonrpc` but doesn't contain json rpc `method` field;
- `invalid_jsonrpc_request` the request contains `jsonrpc` field but it's not a valid json rpc request (cannot decode json body);
- `no_jsonrpc_field` the request does not contain `jsonrpc` field;

To see the errors you can use `LOG_SERVER_DEBUG=true` environment variable. To check full json logs, log in to container and check `/var/log/nginx/access_with_body.json` file.

## Metrics
The interceptor publishes the following metrics:
```
- `ngx_request_count` - the number of requests;
- `ngx_request_size_bytes` - the size of the request;
- `ngx_response_size_bytes` - the size of the response;
- `ngx_request_duration_seconds` - the time of the request;
```
For the `ngx_request_duration_seconds` metric, we use `$request_time` value.
It's the time between the first bytes were read from the client and the log write after the last bytes were sent to the client.


## License
Json RPC interceptor is a Nethermind free and open-source software licensed under the [Apache 2.0 License](https://github.com/NethermindEth/jrpc-interceptor/blob/main/LICENSE) except 3 files that are licensed under the `Mozilla Public License Version 2.0`. See [Notice](https://github.com/NethermindEth/jrpc-interceptor/blob/main/Notice) file for details.