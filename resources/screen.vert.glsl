#version 150

in vec2 pos;
in vec2 tex;

out vec2 vertTexCoord;

void main()
{
    gl_Position = vec4(pos.x, pos.y, 0.0, 1.0); 
    vertTexCoord = tex;
}  
