# Notes

## August 1

It is 10:40pm.
I didn't do much today because of work.
I finished off the .obj loader in the morning.
In the evening I ported it, and the texture loading/watching code, from `learning` into the main source tree.

Need to determine the scope of what I'm doing.
Focus is key.

* 5 floors worth of geometry
   * They can share one megatexture and the same shader
   * No moving parts
   * Need to think of good set pieces for floors 3-5
* 'Crystal' model and shader
* Framebuffer/RTT to upscale 64x64 to a good screen size
* Text splat in the middle of the screen (bitmap font) - needs a shader?
* Crystal collection to unlock progress
* Oh yeah - input processing
* PHYSICS
   * Or at least collision detection and some fakery
   * This is probably the most work and the most unknown

Stretch goals:

* Audio
* Fake HDR (reduce ambient if you're close to light)

What's my goal for tonight?

* Framebuffer
* Have a camera object
* Render a corridor

Game name:

PROJECT CASTLEROOK

After writing this I tried to add a corridor scene, but my OBJ loader code was messed up.
Vertices are fine but normals and texture coords were way off.
Normals seem fixed now but texture coords are still messing up.

## August 2

I'm going back to `learning/mesh` to see if I can fix up the texture coordinates issue in isolation.
I made a bunch of screenshots.
Even a basic cube is messed up!
But planes are okay, so what gives?

According to [the internet](https://gamedev.stackexchange.com/questions/45758/vertex-data-split-into-separate-buffers-or-one-one-structure) I should probably interleave vertex components instead of separating them like I have been.
Maybe doing this will give me an easier time when transforming the OBJ data; though I'd still like to know exactly what's wrong with my existing code.

Ok so it's night and I went back to rework the OBJ loader.
Turns out I was doing the indexing completely wrong.
As in, the final mesh basically ends up with no indexes.
I still have an index and an EBO, but it's just 0 1 2 3 4 ...
Maybe in the future I'll index the OBJ vertices properly.

HUZZAH.
After some fiddling with GIMP and Blender, the corridor scene renders properly!
Now I'm going to implement the 64x64 final resolution so I can get a proper sense of what things will look like.

...

Ok so framebuffering is so close.
I can render scaled up 64x64, but the texture coordinates seem to stop working as soon as I do that.
Like what?

## August 3

The problem is that the UV coordinates are wrong when rendering to the framebuffer.
Added a couple of calls to glGetError but nothing seems to be showing up.

The texture coordinates weren't wrong in the scene render -- they were wrong in the screen quad!
I noticed while uploading a post to Imgur with a comparison that the framebuffered scene didn't just have mirrored UV coordinates, but actual mirrored geometry.
Since it was a rectangular corridor it was a little hard to tell.
But I should have looked at the two renders more carefully.
The fact that the corridor is wobbly actually made it obvious, if I'd flipped back and forth between the two renders.

PROGRESS.
