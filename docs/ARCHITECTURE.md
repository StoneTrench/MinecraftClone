# Architecture

The engine is made of multiple layers, like an onion. Each layer is only aware of the layer below itself. These layer are as follows:
+ Platform
	+ FileSystem
	+ OS calls
+ Core
	+ Logging
	+ Events
	+ Async
	+ Math
+ Resources
	+ File formats
	+ Config
	+ Registry
	+ Language
+ Systems
	+ Modding/WASM
	+ Physics
	+ ECS
	+ Animation
+ Game
	+ Game state
	+ Game tick
	+ Modding api

The following are siblings:
+ Server: synchronises game state between clients.
+ Client: keep track of game state from server. It is the player interface. Rendering.
