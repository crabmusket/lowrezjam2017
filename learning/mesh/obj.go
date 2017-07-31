package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	//"fmt"
)

type Obj struct{
	Name string
	Vertices []float32
	Normals []float32
	TextureCoords []float32
	FaceVerts []uint32
	VertTextureCoords []uint32
	VertNormals []uint32
}

type Object struct{
	Name string
	Positions []float32
	TextureCoords []float32
	Normals []float32
	Indices []uint32
}

type ObjWarning struct {
	File string
	Line int
	Warning string
}

func loadObj(filename string) ([]*Object, []*ObjWarning, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	objs, warnings, err := readObj(file)
	if err != nil {
		return nil, nil, err
	}

	for _, warning := range warnings {
		warning.File = filename
	}

	return objs, warnings, nil
}

func readObj(reader io.Reader) ([]*Object, []*ObjWarning, error) {
	obj := new(Obj)
	var warnings []*ObjWarning

	obj.Name = ""
	index := 0

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		index += 1

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		components := strings.Split(line, " ")
		switch components[0] {
		case "o":
			// TODO
			continue

		case "v":
			if len(components) != 4 {
				warnings = append(warnings, &ObjWarning{
					Line: index,
					Warning: "vertex must have 3 space-separated components",
				})
				continue
			}

			lineWarnings := handleVertex(obj, components[1:])
			for _, warning := range lineWarnings {
				warning.Line = index
			}
			warnings = append(warnings, lineWarnings...)

		case "vn":
			if len(components) != 4 {
				warnings = append(warnings, &ObjWarning{
					Line: index,
					Warning: "vertex normal must have 3 space-separated components",
				})
				continue
			}

			lineWarnings := handleVertexNormal(obj, components[1:])
			for _, warning := range lineWarnings {
				warning.Line = index
			}
			warnings = append(warnings, lineWarnings...)

		case "vt":
			if len(components) != 3 {
				warnings = append(warnings, &ObjWarning{
					Line: index,
					Warning: "texture coordinates must have 2 space-separated components",
				})
				continue
			}

			lineWarnings := handleTextureCoords(obj, components[1:])
			for _, warning := range lineWarnings {
				warning.Line = index
			}
			warnings = append(warnings, lineWarnings...)

		case "f":
			if len(components) != 4 {
				warnings = append(warnings, &ObjWarning{
					Line: index,
					Warning: "face must have 3 space-separated components",
				})
				continue
			}

			lineWarnings := handleFace(obj, components[1:])
			for _, warning := range lineWarnings {
				warning.Line = index
			}
			warnings = append(warnings, lineWarnings...)

		default:
			continue
		}
	}

	return []*Object{obj.Finish()}, warnings, nil
}

func (self Obj) Finish() *Object {
	numVerts := len(self.FaceVerts)

	object := &Object{
		Name: self.Name,
		Indices: make([]uint32, numVerts),
		Positions: make([]float32, len(self.Vertices)),
		TextureCoords: make([]float32, numVerts * 2),
		Normals: make([]float32, numVerts * 3),
	}

	copy(object.Positions, self.Vertices)

	for vert := 0; vert < numVerts; vert += 1 {
		// OBJ indexes are 1-based whereas OpenGL is 0-based
		object.Indices[vert] = self.FaceVerts[vert] - 1

		tObj := vert * 2
		tSelf := (self.VertTextureCoords[vert] - 1) * 2
		copy(object.TextureCoords[tObj:tObj+2], self.TextureCoords[tSelf:tSelf+2])

		nObj := vert * 3
		nSelf := (self.VertNormals[vert] - 1) * 3
		copy(object.Normals[nObj:nObj+3], self.Normals[nSelf:nSelf+3])
	}

	return object
}

func handleVertex(obj *Obj, components []string) []*ObjWarning {
	var warnings []*ObjWarning
	vertex := []float32{0, 0, 0}

	for i := 0; i < 3; i += 1 {
		v, err := strconv.ParseFloat(components[i], 32)
		if err != nil {
			warnings = append(warnings, &ObjWarning{
				Warning: "could not parse vertex component: " + components[i],
			})
		} else {
			vertex[i] = float32(v)
		}
	}

	obj.Vertices = append(obj.Vertices, vertex[0], vertex[1], vertex[2])

	return warnings
}

func handleVertexNormal(obj *Obj, components []string) []*ObjWarning {
	var warnings []*ObjWarning
	normal := []float32{0, 0, 0}

	for i := 0; i < 3; i += 1 {
		v, err := strconv.ParseFloat(components[i], 32)
		if err != nil {
			warnings = append(warnings, &ObjWarning{
				Warning: "could not parse vertex normal: " + components[i],
			})
		} else {
			normal[i] = float32(v)
		}
	}

	obj.Normals = append(obj.Normals, normal[0], normal[1], normal[2])

	return warnings
}

func handleTextureCoords(obj *Obj, components []string) []*ObjWarning {
	var warnings []*ObjWarning
	coords := []float32{0, 0}

	for i := 0; i < 2; i += 1 {
		v, err := strconv.ParseFloat(components[i], 32)
		if err != nil {
			warnings = append(warnings, &ObjWarning{
				Warning: "could not parse texture coordinate: " + components[i],
			})
		} else {
			coords[i] = float32(v)
		}
	}

	obj.TextureCoords = append(obj.TextureCoords, coords[0], coords[1])

	return warnings
}

func handleFace(obj *Obj, components []string) []*ObjWarning {
	var warnings []*ObjWarning
	vertices := []uint32{0, 0, 0}
	texCoords := []uint32{0, 0, 0}
	normals := []uint32{0, 0, 0}
	numVertices := 0
	numTexCoords := 0
	numNormals := 0

	for i := 0; i < 3; i += 1 {
		subcomponents := strings.Split(components[i], "/")

		if len(subcomponents) >= 1 {
			v, err := strconv.ParseInt(subcomponents[0], 10, 32)
			if err != nil {
				warnings = append(warnings, &ObjWarning{
					Warning: "could not parse face vertex index: " + components[i],
				})
			} else {
				vertices[i] = uint32(v)
				numVertices += 1
			}
		}

		if len(subcomponents) >= 2 {
			if len(subcomponents[1]) > 0 {
				v, err := strconv.ParseInt(subcomponents[1], 10, 32)
				if err != nil {
					warnings = append(warnings, &ObjWarning{
						Warning: "could not parse face vertex texture coordinate index: " + components[i],
					})
				} else {
					texCoords[i] = uint32(v)
					numTexCoords += 1
				}
			}
		}

		if len(subcomponents) >= 3 {
			if len(subcomponents[2]) > 0 {
				v, err := strconv.ParseInt(subcomponents[2], 10, 32)
				if err != nil {
					warnings = append(warnings, &ObjWarning{
						Warning: "could not parse face vertex normal index: " + components[i],
					})
				} else {
					normals[i] = uint32(v)
					numNormals += 1
				}
			}
		}

		if len(subcomponents) >= 4 {
			warnings = append(warnings, &ObjWarning{
				Warning: "too many attributes in face element: " + components[i],
			})
		}
	}

	if numVertices == 3 {
		obj.FaceVerts = append(obj.FaceVerts, vertices[0], vertices[1], vertices[2])
	}

	if numTexCoords == 3 {
		obj.VertTextureCoords = append(obj.VertTextureCoords, texCoords[0], texCoords[1], texCoords[2])
	}

	if numNormals == 3 {
		obj.VertNormals = append(obj.VertNormals, normals[0], normals[1], normals[2])
	}

	return warnings
}

func (self Obj) HasVertices() bool {
	return len(self.Vertices) > 0
}

func (self Obj) HasNormals() bool {
	return len(self.Normals) > 0
}

func (self Obj) HasFaceVerts() bool {
	return len(self.FaceVerts) > 0
}

func (self Obj) HasVertTextureCoords() bool {
	return len(self.VertTextureCoords) > 0
}

func (self Obj) HasVertNormals() bool {
	return len(self.VertNormals) > 0
}

func (self Obj) IsEmpty() bool {
	return !self.HasVertices() &&
		!self.HasNormals() &&
		!self.HasFaceVerts() &&
		!self.HasVertTextureCoords() &&
		!self.HasVertNormals()
}
