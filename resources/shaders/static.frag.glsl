#version 150

in vec3 vertNormal;
in vec2 vertTexCoord;

out vec4 colour;

uniform sampler2D textureMap;
uniform float ambient;

void main() {
	float sky = dot(vertNormal, vec3(0.0, 1.0, 0.0)) * 0.25 + 0.75;
	colour = sky * ambient * texture(textureMap, vertTexCoord);
}
