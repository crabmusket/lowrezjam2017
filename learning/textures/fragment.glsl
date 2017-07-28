#version 150

in vec3 vertColour;
in vec2 vertTexCoord;

out vec4 colour;

uniform sampler2D textureMap;

void main() {
	vec4 texColour = texture(textureMap, vertTexCoord);
	if (texColour.a < 0.1) {
		discard;
	}
	colour = texColour * vec4(vertColour, 1.0);
}
