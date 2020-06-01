package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	proto "github.com/gogo/protobuf/proto"
	descriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
)

// Config contains configurations for the plugin
type Config struct {
	Suffix string
	Header *template.Template
	Body   *template.Template
	Footer *template.Template
}

// Plugin implements a protoc generate plugin using a go-template
type Plugin struct {
	config Config
	state  map[string]interface{}
}

// New returns a new plugin for the go template
func New(config Config) *Plugin {
	return &Plugin{config: config}
}

// Generate runs code generator request against the given templates
func (p *Plugin) Generate(in io.Reader, out io.Writer) error {
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("error reading input: %s", err)
	}

	request := new(plugin.CodeGeneratorRequest)
	if err := proto.Unmarshal(input, request); err != nil {
		return fmt.Errorf("error parsing input: %s", err)
	}

	protoFiles := make(map[string]*descriptor.FileDescriptorProto)
	for _, protoFile := range request.GetProtoFile() {
		protoFiles[protoFile.GetName()] = protoFile
	}

	response := new(plugin.CodeGeneratorResponse)

	fileDescriptorSetData, err := proto.Marshal(&descriptor.FileDescriptorSet{File: request.ProtoFile})
	if err != nil {
		return err
	}

	for _, fileName := range request.FileToGenerate {
		if len(fileName) == 0 {
			continue
		}

		protoFile, protoFileFound := protoFiles[fileName]
		if !protoFileFound {
			return fmt.Errorf("%s descriptor not found", fileName)
		}

		values := map[string]interface{}{
			"request":           request,
			"file":              protoFile,
			"fileDescriptorSet": string(fileDescriptorSetData),
		}

		buf := bytes.NewBuffer(nil)

		if p.config.Header != nil {
			if err := p.config.Header.Execute(buf, values); err != nil {
				return err
			}
		}

		if p.config.Body != nil {
			if err := p.config.Body.Execute(buf, values); err != nil {
				return err
			}
		}

		if p.config.Footer != nil {
			if err := p.config.Body.Execute(buf, values); err != nil {
				return err
			}
		}

		filename := strings.Replace(fileName, ".proto", p.config.Suffix, -1)
		bytes := buf.Bytes()

		if strings.HasSuffix(filename, ".go") {
			bytes, err = formatGoSource(bytes)
			if err != nil {
				return err
			}
		}

		content := string(bytes)

		response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
			Name:    &filename,
			Content: &content,
		})
	}

	responseData, err := proto.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshalling output: %s", err)
	}

	if _, err := out.Write(responseData); err != nil {
		return fmt.Errorf("error writing output: %s", err)
	}

	return nil
}

func formatGoSource(in []byte) ([]byte, error) {
	out, err := format.Source(bytes.TrimSpace(in))
	if err != nil {
		return nil, err
	}
	return out, nil
}
