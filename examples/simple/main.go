package main

import (
	b64 "encoding/base64"
	"log"
	"os"
	"text/template"

	descriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	templatePlugin "github.com/orktes/protoc-gen-tools/template"
)

func main() {
	config := templatePlugin.Config{
		Suffix: ".decode.pb.py",
		Body: template.Must(template.New("template.tmpl").Funcs(map[string]interface{}{
			"tf_type": func(val *descriptor.FieldDescriptorProto) string {
				return "tf.float32"
			},
			"tf_name": func(val *descriptor.FieldDescriptorProto) *string {
				return val.Name
			},
			"base64": func(val string) string {
				return b64.StdEncoding.EncodeToString([]byte(val))
			},
		}).ParseFiles("template.tmpl")),
	}
	plug := templatePlugin.New(config)

	if err := plug.Generate(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
