## ssh into the VM

```shell
gcloud compute ssh --zone "europe-west10-a" "data-collector" --project "workflows-413409"
```

## copy files from VM to local (e.g., get the DB)

```shell
gcloud compute copy-files data-collector:/data_collector-main/database.db . --zone=europe-west10-a
```

## send a request to store a measurement in the DB

```shell
curl 34.32.51.179/insert -d '{"tableName": "measurements", "timestamp": 1234, "totalWorkflowDuration": 12.34}'
```