package avro

import "github.com/elodina/syphon/Godeps/_workspace/src/github.com/stealthly/go-avro"

type LogLine struct {
	Line      interface{}
	Source    interface{}
	Tag       map[string]string
	Logtypeid interface{}
	Timings   map[string]int64
}

func NewLogLine() *LogLine {
	return &LogLine{}
}

func (this *LogLine) Schema() avro.Schema {
	if _LogLine_schema_err != nil {
		panic(_LogLine_schema_err)
	}
	return _LogLine_schema
}

// Generated by codegen. Please do not modify.
var _LogLine_schema, _LogLine_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "avro",
    "name": "logLine",
    "fields": [
        {
            "name": "line",
            "default": null,
            "type": [
                "null",
                "string"
            ]
        },
        {
            "name": "source",
            "default": null,
            "type": [
                "null",
                "string"
            ]
        },
        {
            "name": "tag",
            "default": null,
            "type": [
                "null",
                {
                    "type": "map",
                    "values": "string"
                }
            ]
        },
        {
            "name": "logtypeid",
            "default": null,
            "type": [
                "null",
                "long"
            ]
        },
        {
            "name": "timings",
            "default": null,
            "type": [
                "null",
                {
                    "type": "map",
                    "values": "long"
                }
            ]
        }
    ]
}`)
