package gogen

import (
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/utils"
)

func writeProperty(writer io.Writer, name, tag, comment string, tp spec.Type, indent int) error {
	utils.WriteIndent(writer, indent)
	var err error
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
		_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp.Name(), tag, comment)
	} else {
		_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp.Name(), tag)
	}

	return err
}

func getAuths(api *spec.ApiSpec) []string {
	authNames := collection.NewSet()
	for _, g := range api.Service.Groups {
		jwt := g.GetAnnotation("jwt")
		if len(jwt) > 0 {
			authNames.Add(jwt)
		}
	}
	return authNames.KeysStr()
}

func getJwtTrans(api *spec.ApiSpec) []string {
	jwtTransList := collection.NewSet()
	for _, g := range api.Service.Groups {
		jt := g.GetAnnotation(jwtTransKey)
		if len(jt) > 0 {
			jwtTransList.Add(jt)
		}
	}
	return jwtTransList.KeysStr()
}

func getMiddleware(api *spec.ApiSpec) []string {
	result := collection.NewSet()
	for _, g := range api.Service.Groups {
		middleware := g.GetAnnotation("middleware")
		if len(middleware) > 0 {
			for _, item := range strings.Split(middleware, ",") {
				result.Add(strings.TrimSpace(item))
			}
		}
	}

	return result.KeysStr()
}

func responseGoTypeName(r spec.Route, pkg ...string) string {
	if r.ResponseType == nil {
		return ""
	}

	resp := golangExpr(r.ResponseType, pkg...)
	switch r.ResponseType.(type) {
	case spec.DefineStruct:
		if !strings.HasPrefix(resp, "*") {
			return "*" + resp
		}
	}

	return resp
}

func requestGoTypeName(r spec.Route, pkg ...string) string {
	if r.RequestType == nil {
		return ""
	}

	return golangExpr(r.RequestType, pkg...)
}

func golangExpr(ty spec.Type, pkg ...string) string {
	switch v := ty.(type) {
	case spec.PrimitiveType:
		return v.RawName
	case spec.DefineStruct:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("%s.%s", pkg[0], strings.Title(v.RawName))
	case spec.ArrayType:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("[]%s", golangExpr(v.Value, pkg...))
	case spec.MapType:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("map[%s]%s", v.Key, golangExpr(v.Value, pkg...))
	case spec.PointerType:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("*%s", golangExpr(v.Type, pkg...))
	case spec.InterfaceType:
		return v.RawName
	}

	return ""
}
