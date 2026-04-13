# Project HAL

> "I'm sorry Dave, I'm afraid I CAN turn that light off"

A DIY WiFi enabled smart light switch built with a Wemos D1 Mini, 
Go backend, and native iOS app.

# Development Approach
This project was developed with Claude as a collaborative assistant.
Claude was used in the same way a developer might use a knowledgeable 
colleague — for discussing architecture tradeoffs, researching options, 
explaining concepts, and occasional code guidance.

This project intentionally spans less familiar territory including embedded
systems firmware and hardware design. In these areas Claude served a
stronger teaching role, helping build on existing foundational knowledge
and develop genuine understanding of new platforms and concepts.

All architectural decisions were made with genuine understanding of the
tradeoffs involved. Where Claude provided code snippets or suggestions,
they were reviewed and understood before being integrated — with the goal
of learning the platform, not just shipping code.

Code reviews were performed with Claude assistance before each commit,
ensuring consistent quality and catching issues early regardless of
whether the code was written independently or with AI assistance.

This reflects how I believe AI tools should be used in professional 
development — as a force multiplier in familiar territory where existing 
expertise directs the work, and as a way to rapidly build genuine 
understanding in unfamiliar domains. In both cases the developer remains 
in the driver's seat, using AI to move faster and learn deeper rather 
than to bypass understanding altogether.

## Components
- [Firmware](firmware/README.md)
- [Server](server/README.md)
- [iOS App](ios/README.md)
- [Hardware](hardware/README.md)
- [Documentation](docs/README.md)

