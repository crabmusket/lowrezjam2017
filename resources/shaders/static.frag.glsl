#version 150

in vec3 vertNormal;
in vec2 vertTexCoord;
in float vertDist;

out vec4 colour;

uniform sampler2D textureMap;
uniform float ambient;
uniform vec3 fogColour = vec3(0, 0, 0);
uniform float fogStart = 2;
uniform float fogEnd = 10;

void main() {
	float fog = clamp((fogEnd - vertDist) / (fogEnd - fogStart), 0,  1);
	colour = mix(vec4(fogColour, 1), ambient * texture(textureMap, vertTexCoord), fog);
}
