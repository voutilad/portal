# a little portal

Need to expose Google Cloud Storage or Google Cloud Secrets as local files?

Can't use FUSE because your OS/kernel/etc. doesn't support it?

Just use portal!

## Building

`$ make`

Easy enough.

## Running

Currently doesn't leverage environmental discovery of credentials, which
should be a simple change, but for now you need a Google Service Account key
in json form. (Yeah, not the best approach, but this is POC.)

### 1. Open the Portal
Stage the json key and open the portal.

```
$ GOOGLE_APPLICATION_CREDENTIALS=neo4j-se-team-201905-de99ba7a9d55.json \
	./portal --gcs=my-bucket-name,hello.csv:local.csv
Opening portal to GCS { bucket: my-bucket-name, object: hello.csv } @ local.csv
```

### 2. Read from the Portal
In another process or program, just read the file (in this case, `local.csv`).

```
$ cat local.csv
name,age
dave,99
cora,9
maple,6
```

### 3. The Portal will Close
You should see an update from the portal process and it should terminate once
all files are consumed:

```
...
Portal GCS { bucket: my-bucket-name, object: hello.csv } now transmitting...
Wrote 32 bytes from portal GCS { bucket: my-bucket-name, object: hello.csv }
```

## Caveats
1. No environmental auth (yet)
2. Portal closes after all data is consumed...it does not restart or recreate
   the fifo/pipe
3. You must read sequentially, not randomly...this is a pipe, dammit!

