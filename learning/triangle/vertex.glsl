#version 150

in vec3 pos;
in vec3 colour;

out vec3 vertColour;

void main() {
	gl_Position = vec4(pos, 1.0);
	vertColour = colour;
}
