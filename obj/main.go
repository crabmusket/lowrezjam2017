package obj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Material struct{
	Name string
	Start uint32
	End uint32
}

type Object struct{
	Filename string
	Id uint32
	Vertices []float32
	Indices []uint32
	Vbo uint32
	Ebo uint32
	Materials []Material
}

type MaterialData struct{
	Name string
	Start uint32
}

type ObjData struct{
	Vertices []float32
	Normals []float32
	TextureCoords []float32
	FaceVerts []uint32
	VertTextureCoords []uint32
	VertNormals []uint32
	Materials []MaterialData
}

type Warning struct {
	Filename string
	Line int
	Warning string
}

func Load(filename string) (*Object, []*Warning, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	obj, warnings, err := Read(file)
	if err != nil {
		return nil, nil, err
	}

	obj.Bind()

	obj.Filename = filename
	for _, warning := range warnings {
		warning.Filename = filename
	}

	return obj, warnings, nil
}

func Read(reader io.Reader) (*Object, []*Warning, error) {
	obj := new(ObjData)
	var warnings []*Warning

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
				warnings = append(warnings, &Warning{
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
				warnings = append(warnings, &Warning{
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
				warnings = append(warnings, &Warning{
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
				warnings = append(warnings, &Warning{
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

		case "usemtl":
			if len(components) != 2 {
				continue
			}
			handleMaterial(obj, components[1])

		default:
			continue
		}
	}

	object, err := obj.Finish()
	if err != nil {
		return nil, nil, err
	}

	return object, warnings, nil
}

func (self ObjData) Finish() (*Object, error) {
	numVerts := len(self.FaceVerts)
	stride := 3 + 3 + 2

	object := &Object{
		Indices: make([]uint32, numVerts),
		Vertices: make([]float32, numVerts * stride),
		Materials: make([]Material, len(self.Materials)),
	}

	for vert := 0; vert < numVerts; vert += 1 {
		// For now don't bother indexing properly. Just do what you do.
		object.Indices[vert] = uint32(vert)

		begin := vert * stride
		vertIndex := (self.FaceVerts[vert] - 1) * 3
		copy(object.Vertices[begin:begin+3], self.Vertices[vertIndex:vertIndex+3])

		normal := begin + 3
		normalIndex := (self.VertNormals[vert] - 1) * 3
		copy(object.Vertices[normal:normal+3], self.Normals[normalIndex:normalIndex+3])

		texture := normal + 3
		textureIndex := (self.VertTextureCoords[vert] - 1) * 2
		correctCoords := []float32{self.TextureCoords[textureIndex], 1-self.TextureCoords[textureIndex+1]}
		copy(object.Vertices[texture:texture+2], correctCoords)
	}

	if len(object.Vertices) == 0 || len(object.Indices) == 0 {
		return nil, fmt.Errorf("obj data is empty")
	}

	numMaterials := len(self.Materials)
	for i, material := range(self.Materials) {
		object.Materials[i].Name = material.Name
		object.Materials[i].Start = material.Start
		if i < numMaterials-1 {
			object.Materials[i].End = self.Materials[i+1].Start
		} else {
			object.Materials[i].End = uint32(numVerts)
		}
	}

	return object, nil
}

func handleVertex(obj *ObjData, components []string) []*Warning {
	var warnings []*Warning
	vertex := []float32{0, 0, 0}

	for i := 0; i < 3; i += 1 {
		v, err := strconv.ParseFloat(components[i], 32)
		if err != nil {
			warnings = append(warnings, &Warning{
				Warning: "could not parse vertex component: " + components[i],
			})
		} else {
			vertex[i] = float32(v)
		}
	}

	obj.Vertices = append(obj.Vertices, vertex[0], vertex[1], vertex[2])

	return warnings
}

func handleVertexNormal(obj *ObjData, components []string) []*Warning {
	var warnings []*Warning
	normal := []float32{0, 0, 0}

	for i := 0; i < 3; i += 1 {
		v, err := strconv.ParseFloat(components[i], 32)
		if err != nil {
			warnings = append(warnings, &Warning{
				Warning: "could not parse vertex normal: " + components[i],
			})
		} else {
			normal[i] = float32(v)
		}
	}

	obj.Normals = append(obj.Normals, normal[0], normal[1], normal[2])

	return warnings
}

func handleTextureCoords(obj *ObjData, components []string) []*Warning {
	var warnings []*Warning
	coords := []float32{0, 0}

	for i := 0; i < 2; i += 1 {
		v, err := strconv.ParseFloat(components[i], 32)
		if err != nil {
			warnings = append(warnings, &Warning{
				Warning: "could not parse texture coordinate: " + components[i],
			})
		} else {
			coords[i] = float32(v)
		}
	}

	obj.TextureCoords = append(obj.TextureCoords, coords[0], coords[1])

	return warnings
}

func handleFace(obj *ObjData, components []string) []*Warning {
	var warnings []*Warning
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
				warnings = append(warnings, &Warning{
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
					warnings = append(warnings, &Warning{
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
					warnings = append(warnings, &Warning{
						Warning: "could not parse face vertex normal index: " + components[i],
					})
				} else {
					normals[i] = uint32(v)
					numNormals += 1
				}
			}
		}

		if len(subcomponents) >= 4 {
			warnings = append(warnings, &Warning{
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

func handleMaterial(obj *ObjData, name string) {
	obj.Materials = append(obj.Materials, MaterialData{
		Name: name,
		Start: uint32(len(obj.FaceVerts)),
	})
}
