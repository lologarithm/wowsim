# Classic TBC Elemental Shaman DPS Simulator

## Usage

Command line

`go run main.go --config example_config.json`

Important command line options:

`--config`  Location of config file to load. This includes buffs, consumes, gear, gems, enchants, everything about the character

`--rotation`  If you want to test a specific rotation instead of having an AI optimized rotation to maximize mana usage. 
    
  Standard Format:  CL6,LB12,LB12,LB12
    
  Optional 'Priority' casting:   pri,CL6,LB12    (this will cast CL6 anytime off CD, highly likely to go OOM unless fight is short)

If not specified the AI will simply try to use exactly all the mana by casting as many CL as mana will allow.

`--duration`  Number of seconds to run the simulation for. Defaults to 300.

`--iter` Number of iterations to run the simulation for. Defaults to 10,000. Stat weight calculations are more accurate the more iterations run.

`--noopt` No optimizations, disables running gem optimizer and stat weight calculations.


## TODO

### UI
  - Could use a real design
  - Icons
  - Gear Change data on hover of item change dropdown.
  - Remove Gem button

### Items
  - Validate gear stats - unsure if item data source is accurate.
  - More Gems - https://blizzardwatch.com/2021/04/07/burning-crusade-classic-gems/

### Engine
  - Set Bonuses (missing t5/6)

### Other
  - Implement Gear Phases
  - seventyupgrades importer
  - Add armor type to allow for a 'mail-only' optimization
  - 'Gear Sets' both pre-made and let players save the setup. (optionally allow for saving of buffs as well)
  - History - Make another results tab that holds the history of all sims. (probably just Peak DPS + Avg DPS)
  - Versioning - Add a version notification that can do a quick check to see if new version exists. (maybe include a like VERSION file the client can poll on every few minutes)
  - Look into more efficient serialization between sim <-> JS (use same serialization for wasm and server if that is ever implemented)
    - This has been mostly mitigated by sending less data.
  - Write some tests already... so many small breaks from refactors that tests would have caught.

## Install

The simulator and development server can be built and run with Docker.

These commands were tested on Ubuntu 20.04, and may need to be modified for a different OS.

To build (from this directory):

```
docker build -t wowsim .
```

To Run the dev server (from this directory):

```
docker container run -it -p 3333:3333 -v `pwd`/ui:/go/src/app/ui wowsim
```
