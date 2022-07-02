syntax = "v1"

info (
	title: // TODO: add title
	desc: // TODO: add description
	author: "{{.gitUser}}"
	email: "{{.gitEmail}}"
)

type Request {
  Name string `path:"name,options=you|me"`
}

type Response {
  Message string `json:"message"`
}

@server (
    prefix: v1
    group: chat // 分组用来生成更清晰的代码结构，这里强制要求包含
    middleware: test
)
service {{.name}} {
  @handler helloHandler
  get /hello/:name returns (Response)

  @handler echoHandler
  post /echo (Request) returns (Response)
}
