#version 150

in vec2 vertTexCoord;

out vec4 colour;

uniform sampler2D screenTexture;

void main()
{ 
    colour = texture(screenTexture, vertTexCoord);
}
