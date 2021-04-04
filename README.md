# OpenRMS
OpenRMS is an open platform for Slot car Race Management written in Go.

## Motivation
The motivation for starting this project is, my experience getting into
slot cars, and the software tooling around the hobby is for the most part
closed source, and in many cases unmaintained, platform dependent.

The idea is to create an open source race management system capable of
supporting multiple vendors such as Scalextric, Carrera or Oxigen.

Other than supporting as many vendors as possible, The community should
be extend OpenRMS with the features need, it could be better support for
Arduino, it could be having yellow flags light up around the track if a
car has crash.

Last it must be platform independent, the goal is to be able to run it
on a Raspberry Pi, eliminating the high entry price of a computer with
a Windows license. The goal is to release a version that will work on
both Mac OS, Linux and Windows

## Getting Started
### Building
To build run `make openrms`

### Running OpenRMS
When you have open rms build run `./openrms`. if you have want to set
a path for the configuration file set the `-config` flag: `./openrms -config <file>`

## Architecture
OpenRMS is build in a modular way, it's build around 3 different plugin
types:

- Implement
- Rule
- Telemetry

Messages are passed around each component using either events coming
from the `implement`, and commands sent to the `implement`.

Rules are then applied to each event, and each change made by a rule
is sent to the implement via a command. A rule can subscribe to individual
changes, and make changes to the state.

### Implement
The implement is the connector, this plugin type provides connectivity
between your hardware for example the Oxigen Dongle.

### Rule
The rule is the definition of rules which are in effect during a race.
for example fuel simulation

### Telemetry
Optional plugin type which allows all metric collected to be shipped of
to a database like InfluxDB

## Roadmap
- [ ] Web interface
- [ ] GRPC or REST API with streaming data
- [ ] Fuel simulation
- [ ] Damage Simulation
- [ ] Race planning

## Extending

## Contributing

## License
