package api

type Code int

const (
	CodeOK           Code = 101000 + iota // 101000
	CodeUnknown                           // 101002
	CodeInvalidParam                      // 101003
	CodeAPINotFound                       // 101004
)

var _messages = map[Code]string{
	CodeOK:           "ok",
	CodeUnknown:      "unknown error",
	CodeInvalidParam: "invalid paramater",
	CodeAPINotFound:  "API not found",
}

func Message(code Code) string {
	msg, ok := _messages[code]
	if !ok {
		return _messages[CodeUnknown]
	}

	return msg
}
