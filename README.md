Test-Runner
===========

Test-Runner is a Go module which can be used to execute programs against particular
test cases and get corresponding outputs. The goal is to provide an easy module
for testing programs for their correctness. This is one of the many modules
required to create a completely functional backend for an "Online Judge".

## Project Status

The project is in its infancy. A brief roadmap to be completed includes:

- Spawn containers to run jobs.
- Evaluate programs and get output from containers.
- Handle timeouts, errors, etc. Don't let containers die.
- Schedule jobs against available containers.

Other things before release:

- Complete documentation for the project.
- Unit tests for components.
