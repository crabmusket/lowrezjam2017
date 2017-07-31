#version 150

in vec3 pos;
in vec3 norm;
in vec2 tex;

out vec3 vertNormal;
out vec2 vertTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
	gl_Position = projection * view * model * vec4(pos, 1.0);
	vertNormal = (model * vec4(norm, 1.0)).xyz;
	vertTexCoord = tex;
}
