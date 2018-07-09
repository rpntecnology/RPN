package config

const (
	// database server
	DB_SERVER  =  "127.0.0.1:27017"
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