# World generator

+ Initialization
	+ Read all world generators and their manifests.
		+ Create WASM modules.
	+ Parse resources and their types.
	+ Sort phases.
		+ Based on flags and inputs/outputs.

## WorldConfig

+ Each worker gets a reference to a world generator config for that world. Every phase can reference this config and it's fields.
	+ Phase order resolved by the scheduler.
	+ Contains the world seed.
	+ Large structure Registry access.
	+ Small structure Registry access.
	+ Block Registry access.
	+ Entity Registry access.
	+ Loot Registry access.

## Workers

+ Each worker can contain multiple WASM module instances.
+ Wasm modules inside a worker share one memory.
+ Number of workers can be configured.

+ Initialization
	+ Initialize resource locations based on type.
	+ Instantiate WASM modules.
	+ Prepare for region generation.

## Region

+ World is split into regions.
	+ At most one worker per region.
	+ Per region phase cache.
	+ If region is done whole cache can be evicted.

+ World generation is batched, no single chunk generation. But chunks can get priority inside the region.

+ Size of a region is gonna be hard coded as a constant. TODO: tweak size of region for performance.

### Large structures

+ Large structure information (AABB, id) is stored here. Like dungeons, villages.

+ A phase can analyze conditions and emit a structure proposal, which gets added to the region if other conditions apply.
	+ Large structures cannot overlap outside the region for the sake of simplicity.
	+ Large structures shouldn't overlap.

```go

// go:wasmexport phase
func phase() {
	solid := wgen.ReadResource()
	
}

```

## Phase

+ World generation is split into phases, each phase can only read the it's own, and it's 26 (9+8+9) neighbours' previous phase.
+ If there's a cache miss, generate the phase inline.

+ At the final phase, update the region cache and return the function data.

+ Wasm modules get created together, in order of the schedule.
+ Wasm modules get destroyed together with the worker.

### Manifest

```json
{
	"phase": [
		{
			// function name to call, and name to refer to the phase
			"name":"stone",
			"flags": [
				// this phase needs to know about the 26 neighbours
				"neighbours",
				// this phase is the final phase
				"final",
				// this phase is the first phase
				"entry",
			],
			// input resources (using abi calls, used for automatic scheduling and validation)
			"input": [], 
			// output resources (using abi calls, used for automatic scheduling and validation)
			"output": [
				"height [CHUNK_X,CHUNK_Z]i64",
				"solid [CHUNK_X,CHUNK_Y,CHUNK_Z]i64"
			]
		}
	]
}
```

## Noise

+ Noise is sampled and interpolated in the engine after which it crosses the ABI.

## Problems

+ [x] How do large structures know they have to spawn there in the region?
	+ They are proposed by phases.
+ [x] How can multiple mods cooperate and collectively modify the world generation?
	+ Each phase could get a list of tags, or information on their return/input values in the Manifest.
	+ This allows the world generator to interleave phases automatically based on their dependencies.
	+ Or emit warnings when certain phases get left out, or ordering is uncertain/nondeterministic.
+ [x] How do you determine if tho large structures overlap? AABB is simple, but a large tree and a village, shouldn't necesseraly cancel out eachother if their AABBs overlap.
	+ For now AABBs are fine.