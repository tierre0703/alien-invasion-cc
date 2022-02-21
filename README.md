# alien-invasion-cc

Simple Implementation of **Alien Invasion** CLI

---
* [Story](#Story)
* [Implementation](#Implementation)
* [Installation](#Installation)
* [Run](#Run)
* [Test](#Test)
* [Assumption](#Assumption)
---

## Story

Mad​ ​aliens​ ​are​ ​about​ ​to​ ​invade​ ​the​ ​earth​ ​and​ ​you​ ​are​ ​tasked​ ​with​ ​simulating​ ​the invasion.

You​ ​are​ ​given​ ​a​ ​map​ ​containing​ ​the​ ​names​ ​of​ ​cities​ ​in​ ​the​ ​non-existent​ ​world​ ​of X.​ ​The​ ​map​ ​is​ ​in​ ​a​ ​file,​ ​with​ ​one​ ​city​ ​per​ ​line.​ ​The​ ​city​ ​name​ ​is​ ​first,
followed​ ​by​ ​1-4​ ​directions​ ​(north,​ ​south,​ ​east,​ ​or​ ​west).​ ​Each​ ​one​ ​represents​ ​a road​ ​to​ ​another​ ​city​ ​that​ ​lies​ ​in​ ​that​ ​direction.

For​ ​example:

Foo​ ​north=Bar​ ​west=Baz​ ​south=Qu-ux

Bar​ ​south=Foo​ ​west=Bee

The​ ​city​ ​and​ ​each​ ​of​ ​the​ ​pairs​ ​are​ ​separated​ ​by​ ​a​ ​single​ ​space,​ ​and​ ​the
directions​ ​are​ ​separated​ ​from​ ​their​ ​respective​ ​cities​ ​with​ ​an​ ​equals​ ​(=)​ ​sign. You​ ​should​ ​create​ ​N​ ​aliens,​ ​where​ ​N​ ​is​ ​specified​ ​as​ ​a​ ​command-line​ ​argument.

These​ ​aliens​ ​start​ ​out​ ​at​ ​random​ ​places​ ​on​ ​the​ ​map,​ ​and​ ​wander​ ​around​ ​randomly, following​ ​links.​ ​Each​ ​iteration,​ ​the​ ​aliens​ ​can​ ​travel​ ​in​ ​any​ ​of​ ​the​ ​directions
leading​ ​out​ ​of​ ​a​ ​city.​ ​In​ ​our​ ​example​ ​above,​ ​an​ ​alien​ ​that​ ​starts​ ​at​ ​Foo​ ​can​ ​go
north​ ​to​ ​Bar,​ ​west​ ​to​ ​Baz,​ ​or​ ​south​ ​to​ ​Qu-ux.

When​ ​two​ ​aliens​ ​end​ ​up​ ​in​ ​the​ ​same​ ​place,​ ​they​ ​fight,​ ​and​ ​in​ ​the​ ​process​ ​kill each​ ​other​ ​and​ ​destroy​ ​the​ ​city.​ ​When​ ​a​ ​city​ ​is​ ​destroyed,​ ​it​ ​is​ ​removed​ ​from the​ ​map,​ ​and​ ​so​ ​are​ ​any​ ​roads​ ​that​ ​lead​ ​into​ ​or​ ​out​ ​of​ ​it.

In​ ​our​ ​example​ ​above,​ ​if​ ​Bar​ ​were​ ​destroyed​ ​the​ ​map​ ​would​ ​now​ ​be​ ​something like:

Foo​ ​west=Baz​ ​south=Qu-ux

Once​ ​a​ ​city​ ​is​ ​destroyed,​ ​aliens​ ​can​ ​no​ ​longer​ ​travel​ ​to​ ​or​ ​through​ ​it.​ ​This
may​ ​lead​ ​to​ ​aliens​ ​getting​ ​"trapped".

---
## Implementation
* load **cities**, **city links**
* spawn **aliens** in random
* move **aliens** according to **city links**
* if two **aliens** meet in a **city**, do fight:
    * this **city** in fight get destroyed, and removed from **city links**
    * **aliens** in fight will be trapped
* loop the **aliens** move until:
    * maximum **moves** reached
    * all **cities** are destroyed
    * all **aliens** are trapped
---


## Installation
**Step 1: Install Golang** (Go version 1.17+)

**Step 2: Download and Build Source code**
```sh
#Go to source directory
cd alien-invasion-cc

#Build
go build -o bin/alien-invasion-cc

#or
go build
```

**Step 3: Run binary**
```sh
./bin/alien-invasion-cc --help
```
output:
```sh
Usage:
  alien-invasion-cc [flags]

Flags:
  -n, --aliens uint   number of aliens to be spawned (default 5)
  -m, --file string   map file path (default "test_data/test_map")
  -h, --help          help for alien-invasion-cc
  -s, --steps uint    number of maximum moves (default 10000)
```

## Run
```sh
#Run
./bin/alien-invasion-cc -m "test_data/test_map" -n 5 -s 10000
```

The following parameters are available :
* **aliens** (shorthanded to **n**) the number of aliens spawned at startup (defaults to **5**)
* **steps** (shorthanded to **s**) the number of maximum steps allowed (defaults to **10000**)
* **file** (shorthanded to **m**) the path of the world map file (defaults to *test_data/test_map**)

## Test
Run Unit Test
```sh
# Test with code coverage
go test -cover ./...
```

```sh
# Test with verbose output
go test -cover -v ./...
```

## Assumption
1. parameters for **steps** and **aliens** are always positive.
2. **City** names are alpha-numeric only, and no accept for space("space" is reserved for parsing map)
