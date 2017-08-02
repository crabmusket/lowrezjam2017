#version 150

in vec3 vertNormal;
in vec2 vertTexCoord;

out vec3 colour;

uniform sampler2D textureMap;

void main() {
	float d = (dot(vertNormal, vec3(0.0, 1.0, 0.0)) + 1) / 2;
	colour = texture(textureMap, vertTexCoord).xyz * d;
}
