package {{.package}}

import (
    {{.imports}}
)

func (u *{{.service}}) {{.handler}}({{.request}}) {{.response}} {
    // todo: add your logic here and delete this line

	{{.return}}
}