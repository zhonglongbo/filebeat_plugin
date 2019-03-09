package main

// Config for sample processor.
type Config struct {
	Sample         float64        `config:"sample"`     //采样率
	LogType        string         `config:"log_type"`    //采样的日志级别
	QueryType      string         `config:"query_type"`  //采样query类型
}



func defaultConfig() Config {
	return Config{
		Sample: 1,
		LogType: "info",
		QueryType: " ",
	}
}

