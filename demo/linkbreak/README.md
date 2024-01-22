# Link break

A Simple game built from EDD

## Game play

Drag mouse to guide the player object to nearby enemies to chain them. Chained enemies lose overtime then destroyed.

Game over condition: number of enemies appear on screen match or exceed cap.

## System design

- Through automation(system ticks) or user input, actions are perform by `Operators`. `Operators` create and commit `Events` to corresponding `Stores` if the `State transition` is correct.

```
    Automaton/Input -> Operators -> Events -> Stores
```


- The UI displays using `Projection`, created by `Projector`. `Projector` compose `Events` with selective conditions into desire `Projection`.

```
    Stores -> Events -> Projection -> Projector -> UI
```

## Division of State/Entity

States/Entities are divided based on the smallest possible separation of `Events` and `Effects`, Link break separated into the below `Entities`:

- `Health`, contains transitions of health value
- `Position`, contains transitions of position value
- `Runner`, contains data needed to compose a runner (health, position, faction).
- `Link`, contains transition of link value. A `Link` bound between 2 `Runners`(source, and target). 
- `Spawner`, contains data used to spawn a new runner.