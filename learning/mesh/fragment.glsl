#version 150

in vec3 vertNormal;
in vec2 vertTexCoord;

out vec4 colour;

uniform sampler2D textureMap;

void main() {
	colour = texture(textureMap, vertTexCoord) * dot(-vertNormal, vec3(0.0, 1.0, 0.0));
}
