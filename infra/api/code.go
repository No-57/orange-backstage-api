package api

type Code int

const (
	CodeOK           Code = 101000 + iota // 101000
	CodeUnknown                           // 101001
	CodeInvalidParam                      // 101002
	CodeAPINotFound                       // 101003
	CodeForbidden                         // 101004
	CodeNotFound                          // 101005
)

const (
	CodeTokenExpired Code = 102000 + iota // 102000
)

var _messages = map[Code]string{
	CodeOK:           "ok",
	CodeUnknown:      "unknown error",
	CodeInvalidParam: "invalid parameter",
	CodeAPINotFound:  "API not found",
}

func Message(code Code) string {
	msg, ok := _messages[code]
	if !ok {
		return _messages[CodeUnknown]
	}

	return msg
}
