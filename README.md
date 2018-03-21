# poc-go-zap

Playing around with [OWASP ZAP](https://github.com/zaproxy/zaproxy) automation using [zaproxy/zap-api-go](https://github.com/zaproxy/zap-api-go).

## Usage

```
docker-compose up -d
go run *.go
```

The ZAP proxy is available at:

* http://127.0.0.1:8080/

## ZAP Baseline Scan

This will run the baseline scan as configured in `docker-compose-run.yml`:

```
./run-zapbaseline.sh http://example.com/
```

The results are written out to `./reports/`. You can use [jq](https://stedolan.github.io/jq/manual/) to extract various information from the `json` output:

```
jq '.site.alerts[] | "\(.name) \t[\(.riskdesc)]"' ./reports/zap-baseline-example.com.json
```

## Web UI

If you want to use the [ZAP WebSwing UI](https://github.com/zaproxy/zaproxy/wiki/WebSwing), you will have to:

* Change the `zaproxy` service in the ``docker-compose.yml` file to use the `owasp/zap2docker-stable` image
* Change the `zaproxy` command to call `zap-webswing.sh`

Once everything is started up, you can then access the UI at:

* http://127.0.0.1:8080/?anonym=true&app=ZAP

**Note:** It seems that enabling this will break any 'normal' port/proxy capability, including the API. It also seems as though the run script for this doesn't allow command line arguments to be passed to the proxy itself.

## Hackables

* BodgeIt Store: http://127.0.0.1:8081/bodgeit/
* OWASP Juice Shop: http://127.0.0.1:8082/

## Potential Issues

* You can scan the hackables using their 'docker-compose service name' and 'internal port' (as this is from the perspective of the ZAP container), eg.
    * `http://bodgeit:8080/bodgeit/`
    * `http://juiceshop:3000`
* `zaproxy` container logs show error 'URL Not Found in the Scan Tree'
    * You need to access/spider a URL before you can scan it.
    * You may have tried to scan a `127.0.0.1` URL, which is going to reference the ZAP container.. not the local machine.
* `main.go` produces an error such as `spider error: invalid character '<' looking for beginning of value`
    * You're probably running the WebUI version, which seems incompatible with the API..
