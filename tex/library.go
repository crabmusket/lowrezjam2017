package textures

type Library map[string]*Texture

func MakeLibrary() Library {
	return make(Library)
}
