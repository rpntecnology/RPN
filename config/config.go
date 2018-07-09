package config

const (
	// database server
	//DB_SERVER  =  "127.0.0.1:27017"

	// mongodb atlas server
	DB_SERVER = "mongodb://nik0105:xcl940105@cluster0-shard-00-00-tsu2u.gcp.mongodb.net:27017,cluster0-shard-00-01-tsu2u.gcp.mongodb.net:27017,cluster0-shard-00-02-tsu2u.gcp.mongodb.net:27017/test?replicaSet=Cluster0-shard-0&authSource=admin"

	// database name
	DB_NAME =  "RPN_proj"

	BUCKET_NAME = "post-images-2039211"

)

var MySigningKey = []byte("mySecret")

var (
	MediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
		".gif":  "image",
		".png":  "image",
	}
)