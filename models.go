package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeType string;

const (
	FileType NodeType = "file"
	DirectoryType NodeType = "directory"
)


type FileSystemNode struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Path string `bson:"path"`
	FileName string   `bson:"name"`
	Size   int64    `bson:"size"`
	CID   string   `bson:"CID"`
}


