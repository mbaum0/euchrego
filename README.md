# Euchre in Go (euchrego)

## Overview
My day job frequently requires me to write in Go but I often feel that don't know the more subtle quirks of the language. I wanted to work on a project from the ground up that pushes me to experiment with the language.

I am a big fan of the card Game Euchre and decided to make a program that would allow me to play remotely with my family and friends. This would hopefully help me touch on the follow concepts in Go:
 - Game design
 - State machines
 - Advanced terminal output (first iteration will include a PTUI)
 - Networking (for multiplayer)
 - Database interations (for game state)
 - Event-driven FSM (allows recreation of state from a sequence of events)

There are probably other things I'll end up diving into as I work through this.

## Design
I've decided to use an FSM to tackle game state and progression as I've designed them in the past in C. FSMs are a concept I'm comfortable with and seem like the perfect tool for the job.

I believe it will also enable me to create a predictable system where I can recreate state from a saved list of transitions for easy unit testing, while also being able to save games and pick up where they left off. 

The state flow diagram is included in this repo under `fsm.md`. It is written in mermaid and GH should be able to render it.