package analysis

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	dart "github.com/UserNobody14/tree-sitter-dart/bindings/go"
	sitter "github.com/smacker/go-tree-sitter"
)

type Field struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsStatic bool   `json:"isStatic"`
	IsFinal  bool   `json:"isFinal"`
	IsConst  bool   `json:"isConst"`
}

type Method struct {
	Name       string            `json:"name"`
	ReturnType string            `json:"returnType"`
	IsStatic   bool              `json:"isStatic"`
	IsAbstract bool              `json:"isAbstract"`
	Parameters []MethodParameter `json:"parameters"`
}

type MethodParameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type DartInterface struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Superclass string   `json:"superclass,omitempty"`
	Interfaces []string `json:"interfaces,omitempty"`
	Fields     []Field  `json:"fields"`
	Methods    []Method `json:"methods"`
}

func isPublic(name string) bool {
	return !strings.HasPrefix(name, "_")
}

// Parse a single .dart file and extract public interfaces
func parseDartFile(filename string, src []byte) ([]DartInterface, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(dart.GetLanguage())

	tree := parser.Parse(nil, src)
	root := tree.RootNode()

	var interfaces []DartInterface

	visit := func(node *sitter.Node) {
		// Parse classes, mixins, extensions, typedefs
		switch node.Type() {
		case "class_declaration", "mixin_declaration", "extension_declaration", "type_alias_declaration":
			var i DartInterface
			i.Type = map[string]string{
				"class_declaration":      "class",
				"mixin_declaration":      "mixin",
				"extension_declaration":  "extension",
				"type_alias_declaration": "typedef",
			}[node.Type()]

			// Name
			for i := 0; i < int(node.NamedChildCount()); i++ {
				child := node.NamedChild(i)
				if child.Type() == "identifier" {
					i.Name = child.Content([]byte(src))
					break
				}
			}
			if !isPublic(i.Name) {
				return // Skip private
			}

			// Superclass / interfaces (for class/mixin)
			if i.Type == "class" || i.Type == "mixin" {
				for j := 0; j < int(node.NamedChildCount()); j++ {
					child := node.NamedChild(j)
					if child.Type() == "superclass" {
						i.Superclass = child.NamedChild(0).Content([]byte(src))
					}
					if child.Type() == "interfaces" {
						for k := 0; k < int(child.NamedChildCount()); k++ {
							i.Interfaces = append(i.Interfaces, child.NamedChild(k).Content([]byte(src)))
						}
					}
				}
			}
			// Fields and methods
			for j := 0; j < int(node.NamedChildCount()); j++ {
				child := node.NamedChild(j)
				if child.Type() == "class_body" || child.Type() == "mixin_body" || child.Type() == "extension_body" {
					for k := 0; k < int(child.NamedChildCount()); k++ {
						member := child.NamedChild(k)
						switch member.Type() {
						case "method_declaration", "getter_signature", "setter_signature":
							method := Method{
								Name:       extractMethodName(member, src),
								ReturnType: extractReturnType(member, src),
								IsStatic:   hasChildOfType(member, "static_modifier"),
								IsAbstract: hasChildOfType(member, "abstract_modifier"),
							}
							if isPublic(method.Name) {
								method.Parameters = extractParameters(member, src)
								i.Methods = append(i.Methods, method)
							}
						case "field_declaration":
							for l := 0; l < int(member.NamedChildCount()); l++ {
								field := member.NamedChild(l)
								if field.Type() == "variable_declarator" {
									f := Field{
										Name:     field.NamedChild(0).Content([]byte(src)),
										Type:     extractFieldType(member, src),
										IsStatic: hasChildOfType(member, "static_modifier"),
										IsFinal:  hasChildOfType(member, "final_modifier"),
										IsConst:  hasChildOfType(member, "const_modifier"),
									}
									if isPublic(f.Name) {
										i.Fields = append(i.Fields, f)
									}
								}
							}
						}
					}
				}
			}
			interfaces = append(interfaces, i)
		}
		// Visit children recursively
		for i := 0; i < int(node.NamedChildCount()); i++ {
			visit(node.NamedChild(i))
		}
	}
	visit(root)
	return interfaces, nil
}

func extractMethodName(node *sitter.Node, src []byte) string {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == "identifier" {
			return child.Content(src)
		}
	}
	return ""
}

func extractReturnType(node *sitter.Node, src []byte) string {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == "type" {
			return child.Content(src)
		}
	}
	return "void"
}

func hasChildOfType(node *sitter.Node, typ string) bool {
	for i := 0; i < int(node.ChildCount()); i++ {
		if node.Child(i).Type() == typ {
			return true
		}
	}
	return false
}

func extractParameters(node *sitter.Node, src []byte) []MethodParameter {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == "formal_parameter_list" {
			var params []MethodParameter
			for j := 0; j < int(child.NamedChildCount()); j++ {
				p := child.NamedChild(j)
				if p.Type() == "required_parameter" || p.Type() == "default_parameter" {
					paramType := ""
					name := ""
					for k := 0; k < int(p.NamedChildCount()); k++ {
						if p.NamedChild(k).Type() == "type" {
							paramType = p.NamedChild(k).Content(src)
						}
						if p.NamedChild(k).Type() == "identifier" {
							name = p.NamedChild(k).Content(src)
						}
					}
					params = append(params, MethodParameter{
						Name:     name,
						Type:     paramType,
						Required: p.Type() == "required_parameter",
					})
				}
			}
			return params
		}
	}
	return nil
}

func extractFieldType(node *sitter.Node, src []byte) string {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == "type" {
			return child.Content(src)
		}
	}
	return ""
}

// Scan all .dart files in a directory recursively
func ExtractDartInterfaces(dir string) ([]DartInterface, error) {
	var interfaces []DartInterface
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".dart") {
			src, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			ifaces, err := parseDartFile(path, src)
			if err != nil {
				return err
			}
			interfaces = append(interfaces, ifaces...)
		}
		return nil
	})
	return interfaces, err
}

// Gera o arquivo interfaces.json no diretório dado
func GenerateInterfacesJSON(srcDir, outPath string) error {
	ifaces, err := ExtractDartInterfaces(srcDir)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(ifaces, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outPath, data, 0644)
}
