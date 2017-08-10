#version 150

in vec3 vertPos;
in vec3 vertNormal;
in vec2 vertTexCoord;
in float vertDist;

out vec4 fragColour;

#define POINT_LIGHT_COUNT 6
struct PointLight {
	vec3 position;
	vec3 diffuseColour;
	float radius;
};

uniform sampler2D textureMap;
uniform vec3 cameraPos;
uniform float ambient;
uniform vec3 fogColour = vec3(0, 0, 0);
uniform float fogStart = 2;
uniform float fogEnd = 10;
uniform PointLight pointLights[POINT_LIGHT_COUNT];

vec3 handlePointLight(PointLight light, vec3 pos, vec3 normal, vec3 viewDir);

void main() {
	vec3 normal = normalize(vertNormal);
	vec3 viewDir = normalize(cameraPos - vertPos);

	vec3 light = vec3(ambient, ambient, ambient);
	for (int i = 0; i < POINT_LIGHT_COUNT; i += 1) {
		light += handlePointLight(pointLights[i], vertPos, normal, viewDir);
	}

	float fog = clamp((fogEnd - vertDist) / (fogEnd - fogStart), 0,  1);

	vec4 textureColour = texture(textureMap, vertTexCoord);
	fragColour = mix(vec4(fogColour, 1), vec4(light, 1) * textureColour, fog);
}

vec3 handlePointLight(PointLight light, vec3 pos, vec3 normal, vec3 viewDir) {
	vec3 lightDir = normalize(light.position - pos);
	float diffuse = max(dot(normal, lightDir), 0);
	float distance = length(light.position - pos);
	float attenuation = clamp((light.radius - distance) / light.radius, 0, 1);
	return attenuation * light.diffuseColour;
}
