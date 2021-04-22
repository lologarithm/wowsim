# Classic TBC Elemental Shaman DPS Simulator

## TODO

### UI
  - Could use a real design
  - Icons
  - Gear Quality Colors
  - Gear Change data on hover of item change dropdown.
  - Remove Gem button

### Items
  - Gear Quality
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
  - Pull out advanced controls pane - 
    - 'Remove All Gear'
    - 'weights value modulation'
    - 'Gem Everything'
    - '70Upgrades Import'
    - 'Phase Selector'
    - 'Hide Green/Blue Items'
  - Look into more efficient serialization between sim <-> JS (use same serialization for wasm and server if that is ever implemented)
  
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
