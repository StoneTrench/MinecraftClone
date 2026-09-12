# Here's the plan

Ship and Active level mean the same thing.

# Background systems

## Networking


## Physics
- [ ] Collisions. Entities are mapped to chunks, so they are essentially spatially mapped, only try to collide entities which are in the same chunk, or neighbouring chunks.
	- [ ] Entities have simple either spherical, or cylindrical, or maybe even cuboidal axis aligned collision boxes.
	- [ ] Ships have a large AABB, which wraps around the whole thing, to quickly determine if there could be a collision, then their collision has to be handled per voxel, since they can be concave, maybe it could be further refined by having a non axis aligned bounding box too.

- [ ] Decouple physics and fps, this mainly applies to the client, where if the fps is high we should accumulate physics ticks, and vice versa.
	- [ ] Physics happens at a fixed `delta_time`.

## Modding
+ [x] Load mod headers
+ [x] Sort based on dependencies
+ [x] Skip mods with missing dependencies, and collect errors
+ [ ] Load mod assets
+ [x] Load mod program
+ [x] Load api for mod
+ [x] Call initializer functions (to define registry entries, callbacks, systems)

- [ ] World generation happens in phases, where each phase only has access to the previous phase data in itself and neighbouring chunks.
	- [ ] World generation doesn't read data from the actual world, but either from a world generation cache, or a it regenerates the data for the neighbours.
	- [ ] This is so that multiple chunks can be generated in parallel, without them having to share data (except for maybe read from a cache, but only read from it)

- [ ] Implement configuration library, which handles engine configuration (physics tick speed, and other stuff I can't think of right now), and can also be used by the mods, to store configuration.

## Other
- [ ] Entities are mapped to chunks
- [ ] Ships can be in multiple chunks at the same time (decided by their AABB)

- [ ] Certain systems should only run if there's time for it.
	- [ ] Essentially there could be a priority queue for tasks the engine needs to do, where high priority behaviour, like rendering and input handling shouldn't be blocked by things like world generation and chunk mesh generation etc.

- [ ] Event system

# Rendering
- [ ] Create pipeline helper functions

- [ ] Configuration at runtime 
	- [ ] Switch GPU
	- [ ] GUI scale
	- [ ] Set window mode
	- [ ] Reload shaders
	- [ ] Reload textures

- [ ] Global camera.
- [ ] Skybox.
- [ ] Level rendering.
	- [ ] Ship rendering.

## Resources
- [ ] WebAssembly target compatible file system wrapper
	- So we can embed the engine+mods into a webpage
- [ ] Unified asset loading from mods
- [ ] (client side only) Texture Atlas, texture loading from mods
- [ ] (client side only) Sounds registry

# Features
- [ ] Insert items and icons into text like in Factorio.
- [ ] Blocks can have clickable/interactable boxes (like in create) on their surface.
- [ ] Recipes on machine can be easily configured (Note: I don't know yet what machines will be, how they will interact, and what they will do (probably lots of multiblock structures)).
- [ ] When you hold LeftAlt you can move the mouse around like a point and click adventure.
- [ ] No ores I guess, you gotta process different rock types and get ores that way. But what is the point of mining then? (build excavators I guess)

- [ ] Ability to reorganize creative item catalogue in game, also sorting the inventory
- [ ] Every item should have a large amount of tags, so that mods can reason about them

- [ ] Computers, lots of computers, fast computers, using some sort of machine code, that the player can upload to them.
	- [ ] Lots of sensors, for the computers, so they can detect, and analyze their surroundings, and communicate with eachother.

- [ ] Ships should be allowed to be big, huge even! (I don't entirely know how I would implement that, I was even considering making the static level an entity too, but what would the point of that be? When it would be the only entity like that in the entire ECS.)

- [ ] Robotic arm, which can interact with the world the same way a player can, but can be programmed to do things (pick and place, click, break, weld, activate etc)
	- [ ] One kind of arm can interact with ships (hold onto them and try to move them if it's not too heavy) this is gonna be large and expensive.
	- [ ] And one kind that is small, and is more or less handled like a simple player entity, this one is smaller and slightly cheaper.

- [ ] There needs to be a lot of mechanical implements, they have to be easily configurable, and they should be driveable mechanically, but also configured to move freely.
	- [ ] Rotor
	- [ ] Hinge
	- [ ] Ball joint
	- [ ] Linear bearing
	- [ ] Radial rotor (like a tank turrets ring, has to be build in more or less a circular shape (again, multiblock structure))
	- [ ] Piston

- [ ] Fuck you and your crafting grid, terraria does it better, but also, should you be able to craft with your bare hands? Or should it require automation? (Obviously without making it tedious)
	- [ ] Maybe you can craft things by hand, but it's more expensive/slow, except for a couple very expensive things, which actually need production lines.

- [ ] Should there be a tech tree? Or is tech just locked behind you putting together the production/resource extraction to get that thing you want.

- [ ] Also something like the schematic-cannon in create is a must have, but done differently so it's not jarring to see the same exact ideas stolen.
	- [ ] I just can't think of a better idea rightnow. (maybe It could work like in space engineers)
	- [ ] Or maybe construction drones could build it for you like in factorio.