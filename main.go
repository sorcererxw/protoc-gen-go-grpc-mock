package main

import (
	_ "embed"
	"fmt"
	pathpkg "path"
	"sort"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type Data struct {
	Source   string
	Package  string
	Imports  []string
	Services []Service
}

type Service struct {
	Service string
	Funcs   []Func
}

type Func struct {
	Func     string
	Request  string
	Response string
}

func parseFile(path string, file *protogen.File) (data Data) {
	data.Source = path
	data.Package = string(file.GoPackageName)

	imports := make(map[string]string)

	getMessageName := func(m *protogen.Message) string {
		importPath := string(m.GoIdent.GoImportPath)
		if importPath == `./` {
			return m.GoIdent.GoName
		}
		alias, ok := imports[importPath]
		if !ok {
			alias = pathpkg.Base(importPath)
			imports[importPath] = alias
		}
		return alias + "." + m.GoIdent.GoName
	}

	for _, s := range file.Services {
		svc := Service{
			Service: s.GoName,
		}
		for _, m := range s.Methods {
			svc.Funcs = append(svc.Funcs, Func{
				Func:     m.GoName,
				Request:  getMessageName(m.Input),
				Response: getMessageName(m.Output),
			})
		}
		data.Services = append(data.Services, svc)
	}
	for k, v := range imports {
		data.Imports = append(data.Imports, fmt.Sprintf(`%s "%s"`, v, k))
	}

	return orderData(data)
}

func orderData(data Data) Data {
	services := data.Services
	sort.Slice(services, func(i, j int) bool {
		return services[i].Service < services[j].Service
	})
	for k := range services {
		funcs := services[k].Funcs
		sort.Slice(funcs, func(i, j int) bool {
			return funcs[i].Func < funcs[j].Func
		})
		services[k].Funcs = funcs
	}
	data.Services = services
	return data
}

//go:embed mock.tmpl
var _MockTmpl string

var tpl = template.Must(template.New("").Parse(_MockTmpl))

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for path, file := range plugin.FilesByPath {
			if !file.Generate {
				continue
			}
			data := parseFile(path, file)
			if len(data.Services) == 0 {
				continue
			}
			return tpl.Execute(
				plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_grpc_mock.pb.go", file.GoImportPath),
				data,
			)
		}
		return nil
	})
}
