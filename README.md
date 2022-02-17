# a little portal

Need to expose Google Cloud Storage or Google Cloud Secrets as local files?

Can't use FUSE because your OS/kernel/etc. doesn't support it?

Just use portal!

## Building

You need Go 1.16 or newer. Then:

`$ make`

Easy enough?

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

### 3. The Portal will Recycle until you Kill It
You should see an update from the portal process:

```
...
Portal GCS { bucket: my-bucket-name, object: hello.csv } now transmitting...
Wrote 32 bytes from portal GCS { bucket: my-bucket-name, object: hello.csv }
```

Use `ctrl-c` or whatnot to interrupt and kill it :-)


## Caveats
1. No environmental auth (yet!)
2. You must read sequentially, not randomly...this is a pipe, dammit!
3. If you're using with `neo4j-admin import`, because of something (Java?) doing
   a dance to check for magic bytes indicating GZip or Zip files, it will most
   likely drop the first 4 bytes of each file. I recommend:

   a. Use a single file that includes the header you want.
   b. Make sure the first 4 characters can be possible throwaway.

For example:

```
$ head -n2 random-nodes.csv
junkid:ID,name,age:int,:LABEL
1,Person 1,1,FakePerson

$ head -n2 random-rels.csv
crap:START_ID,weight:float,:END_ID,:TYPE
10420,0.184856666,10594,FOLLOWS
```
