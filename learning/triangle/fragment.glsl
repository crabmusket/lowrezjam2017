#version 150

in vec3 vertColour;

out vec4 colour;

void main() {
	colour = vec4(vertColour, 1);
}
