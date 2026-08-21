 # CC:Robot Project

An automation platform for [CC:Tweaked](https://tweaked.cc/) in Minecraft.
The project will provide an API bridge for Lua clients, a Go command-line
interface for controlling those clients, and a small scripting language for
describing in-game automation tasks.

## Goals

- Connect to and manage multiple CC:Tweaked Lua clients.
- Expose a stable API for client communication and control.
- Provide a Go CLI for sending commands, monitoring clients, and running jobs.
- Create a simple scripting language for reusable automation workflows.
- Make task execution observable, scriptable, and safe to stop.

## Planned Architecture

```text
Lua client(s) <──> API bridge <──> Go CLI
								  └──> automation scripts
```

The bridge will handle client connections and message routing. The CLI will
provide operational control, while the scripting layer will translate human-
readable workflows into commands executed by Lua clients in Minecraft.

## Example Workflow

```text
connect robot-1
run mining_job.script
status robot-1
stop robot-1
```

## Status

Early planning and design.

## Contributing

Design discussions, implementation ideas, and CC:Tweaked compatibility notes
are welcome as the architecture develops.
