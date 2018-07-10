package config

const (
	// database server
	//DB_SERVER  =  "127.0.0.1:27017"

	// mongodb atlas server 3.4 or earlier
	DB_SERVER = "mongodb://nik0105:xcl940105@cluster0-shard-00-00-tsu2u.gcp.mongodb.net:27017,cluster0-shard-00-01-tsu2u.gcp.mongodb.net:27017,cluster0-shard-00-02-tsu2u.gcp.mongodb.net:27017/test?replicaSet=Cluster0-shard-0&authSource=admin"

	// mongodb atlas server 3.6
	//DB_SERVER = "mongodb+srv://nik0105:xcl940105@cluster0-tsu2u.gcp.mongodb.net/test"

	// database name
	DB_NAME =  "RPN_proj"

	BUCKET_NAME = "post-images-rpnserver"

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