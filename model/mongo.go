package model

type Trainer struct {
	Content string `bson:"content"` // bson mongo 支持的 tag
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
	Read uint `bson:"read"` // 是否已读，已读为1
}

type Result struct {
	StartTime int64
	Msg string
	Content interface{}
	From string
}
