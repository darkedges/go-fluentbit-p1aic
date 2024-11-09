# go-fluentbit-p1aic

This is a FluentBit input module for the ingestion of PingOne Advanced Identity Cloud logs. It is a clone of the capability developed by [Jon Knight](https://splunkbase.splunk.com/apps?author=jonkenator) for his Splunk App [PingOne Advanced Identity Cloud App](https://splunkbase.splunk.com/app/7529)

## Configuration

| key                      | value                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `p1aic_id_cloud_tenant ` | Enter the Identity Cloud tenant URL, for example, `https://<tenant-env-fqdn>`, where `<tenant-env-fqdn>` is your Identity Cloud tenant.                                                                                                                                                                                                                                                                                                                                                 |
| `api_key_id`             | Enter the API Key ID that you'll use to authenticate to the Identity Cloud REST API endpoints.   See <https://backstage.forgerock.com/docs/idcloud/latest/developer-docs/authenticate-to-rest-api-with-api-key-and-secret.html#get_an_api_key_and_secret>                                                                                                                                                                                                                               |
| `api_key_secret`         | Enter the API Key Secret that you'll use to authenticate to the Identity Cloud REST API endpoints.                                                                                                                                                                                                                                                                                                                                                                                      |
| `log_sources`            | Provides a comma-separated list of Identity Cloud log sources to capture. By default, the app `captures am-authentication`, `am-access`, `am-config` and `idm-activity`. You can add more log sources if needed, but be aware that the data returned by some logs may not be searchable. See Source descriptions for further information on Identity Cloud logging sources. See <https://backstage.forgerock.com/docs/idcloud/latest/tenants/audit-debug-logs.html#source-descriptions> |
| `log_filter`             | Filter to be applied. Example `/payload co "WARNING"`. See <https://backstage.forgerock.com/docs/idcloud/latest/tenants/audit-debug-logs.html#filter-log-results>                                                                                                                                                                                                                                                                                                                       |
| `db`                     | location od state file. example `/fluentbit/state.json`                                                                                                                                                                                                                                                                                                                                                                                                                                 |

## Playground                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |


Edit <fluent-bit.conf.example> and rename to `fluent-bit.conf`

Build the docker image locally to see how it works.

```bash
docker build . -t fluent-bit-p1aic
ocker run -it --rm -v ${PWD}/state:/fluentbit fluent-bit-p1aic 
```

The output produced should resemble the following:
```
Fluent Bit v3.1.10
* Copyright (C) 2015-2024 The Fluent Bit Authors
* Fluent Bit is a CNCF sub-project under the umbrella of Fluentd
* https://fluentbit.io

______ _                  _    ______ _ _           _____  __  
|  ___| |                | |   | ___ (_) |         |____ |/  | 
| |_  | |_   _  ___ _ __ | |_  | |_/ /_| |_  __   __   / /`| | 
|  _| | | | | |/ _ \ '_ \| __| | ___ \ | __| \ \ / /   \ \ | | 
| |   | | |_| |  __/ | | | |_  | |_/ / | |_   \ V /.___/ /_| |_
\_|   |_|\__,_|\___|_| |_|\__| \____/|_|\__|   \_/ \____(_)___/

[2024/11/07 05:57:11] [ info] [fluent bit] version=3.1.10, commit=e28f447995, pid=1
[2024/11/07 05:57:11] [ info] [storage] ver=1.5.2, type=memory, sync=normal, checksum=off, max_chunks_up=128
[2024/11/07 05:57:11] [ info] [cmetrics] version=0.9.7
[2024/11/07 05:57:11] [ info] [ctraces ] version=0.5.6
[2024/11/07 05:57:11] [ info] [input:p1aic:p1aic.0] initializing
[2024/11/07 05:57:11] [ info] [input:p1aic:p1aic.0] storage_strategy='memory' (memory only)
No previous beginTime saved so backdating to 1 minute ago
[2024/11/07 05:57:11] [ info] [input:p1aic:p1aic.0] thread instance initialized
[2024/11/07 05:57:11] [ info] [output:stdout:stdout.0] worker #0 started
[2024/11/07 05:57:12] [ info] [output:azure_blob:azure_blob.1] account_name=darkedges, container_name=logs, blob_type=appendblob, emulator_mode=no, endpoint=darkedges.blob.core.windows.net, auth_type=key
[2024/11/07 05:57:12] [ info] [sp] stream processor started
[0] p1aic: [[1730958939.187197859, {}], {"http"=>{"request"=>{"secure"=>true, "headers"=>{"content-type"=>["application/x-www-form-urlencoded"], "host"=>["am.fr-platform"], "user-agent"=>["Go-http-client/1.1"]}, "method"=>"POST", "path"=>"https://am.fr-platform/am/oauth2/introspect"}}, "level"=>"INFO", "realm"=>"/", "timestamp"=>"2024-11-07T05:55:39.186Z", "topic"=>"access", "transactionId"=>"a431133c-ee0f-46ab-92fc-a9802ce43485/0/0", "component"=>"OAuth", "eventName"=>"AM-ACCESS-ATTEMPT", "request"=>{"detail"=>{}}, "source"=>"audit", "_id"=>"970c772a-f46d-4620-9884-763c697c90fd-44956", "client"=>{"ip"=>"10.100.0.238", "port"=>44844.000000}}]
```