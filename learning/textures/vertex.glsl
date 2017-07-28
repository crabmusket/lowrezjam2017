#version 150

in vec3 pos;
in vec3 colour;
in vec2 tex;

out vec3 vertColour;
out vec2 vertTexCoord;

void main() {
	gl_Position = vec4(pos, 1.0);
	vertColour = colour;
	vertTexCoord = tex;
}
