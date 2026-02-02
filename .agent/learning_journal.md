# Learning Journal

This document tracks the "Why" and "How" of our engine development, focusing on Go concepts and architectural patterns for an intermediate programmer.

---

## 1. Project Foundations: Why Go and Ebiten?

### Why Go?
For a game engine, we need performance and control, but as a solo/small team, we also need productivity.

*   **Composition over Inheritance**: Go does not have `class`, `extends`, or inheritance. It uses **structs** and **interfaces**.
    *   *Traditional OOP*: You might have `GameObject -> MovingObject -> Character -> Player`. This creates "fragile base classes" where changing the parent breaks the child.
    *   *Go's Way*: We define small, independent pieces (Components) and glue them together. This naturally leads to **Entity-Component-System (ECS)** usage, which is the gold standard for game engines.
*   **Performance vs Complexity**: Go offers performance close to C++ but with the readability of Python/JavaScript. The Garbage Collector (GC) handles memory for us, though we must be careful not to overwork it (generating trash) in the game loop (60fps).
*   **Readability**: Go code is famous for being "boring." There is no magic. If code runs, you usually see it explicitly called. This makes debugging an engine much easier.

### Why Ebiten?
Ebiten is a library, not a framework.

*   **The "Just Draw It" Philosophy**: Ebiten gives us a `Draw(screen)` function called 60 times a second. It gives us an empty canvas. It does **not** give us:
    *   Physics
    *   Scene management
    *   Entity lists
    *   Collisions
*   **Why this is good for learning**: Because Ebiten doesn't provide these, **we have to build them.**
    *   If we used Godot, we would just "configure" a RigidBody.
    *   With Ebiten, we have to *write* the code that checks if two rectangles overlap.
    *   This forces you to understand the *architecture* of a game engine from the ground up.

### Key Takeaway
We are building the **logic** (the "Engine") on top of Ebiten's **rendering** (the "Graphics").

---

## 2. Memory Management in Go Games
*From `best_practices.md`*

Go has a Garbage Collector (GC), which automatically frees memory. This is great for web servers, but tricky for games.

*   **The Problem**: If you create a new struct (e.g., `new Bullet()`) 60 times a second, the GC eventually has to "stop the world" to clean up thousands of dead bullets. This causes a "lag spike" (hiccup) in the game.
*   **The Solution**: **Pre-allocation** and **Pooling**.
    *   Instead of `new Bullet()`, we make a slice `bullets := make([]Bullet, 1000)` at the start.
    *   When we need a bullet, we grab an unused one from the array.
    *   When it dies, we just mark it as "unused" (or dead). We never actually "delete" it from memory until the game ends.
    *   **Result**: Zero memory allocations during gameplay = Zero GC pauses = Smooth 60 FPS.

### The Trade-off: RAM vs CPU
*   **This uses more RAM.** If you pre-allocate 10,000 bullets but only use 5, you are "wasting" memory for 9,995 bullets.
*   **Why we do it**: RAM is cheap and abundant (GBs). CPU time in a game loop is scarce (16 milliseconds per frame).
*   **Golden Rule**: We trade *space* (RAM) to gain *time* (Speed).

---

## 3. Profiling: The "X-Ray" for Code
*From `best_practices.md`*

You can't fix what you can't see. "Lag" is invisible. Is it the rendering? The physics? The AI?

*   **`go tool pprof`** is Go's built-in detective.
*   **How it works**: It interrupts your program 100 times a second and writes down "What function is running right now?"
*   **The Result**: A graph (flame graph) showing exactly where your game spends its time.
    *   *Example*: You might think your AI is slow, but `pprof` shows you are actually spending 90% of your time creating new strings for the score display.
*   **Why use it**: Optimizing without profiling is just guessing. Always measure first.

---

## 4. Concurrency: Goroutines & Channels
*From `best_practices.md`*

This is Go's "Superpower."

### The Problem
In a game, the **Main Loop** (Update/Draw) must run every 16ms. If you stop to download a file or wait for a network packet, the game freezes.

### The Solution
*   **Goroutine (`go func()`)**: A "lightweight thread." It's like hiring a background worker.
    *   *Code*: `go DownloadFile("level2.png")`
    *   *Effect*: The main game loop continues running smoothly while the file downloads in the background.
*   **Channel (`chan`)**: A pipe to talk to that worker.
    *   *Analogy*: You don't shout at the worker. You send a note through a pneumatic tube.
    *   *Usage*: The worker finishes downloading and sends the image back through the channel `results <- image`. The Main Loop checks the mailbox `img := <- results` and displays it.

**Why "Idiomatic"?**
In other languages (C++, Java), threads share memory, which leads to crashes (Race Conditions). In Go, we **"Share memory by communicating"** (using channels). Because the worker *sends* the data to you, you don't fight over it.

---

## 5. The First Breath: Wiring the Engine
*What we built in `main.go`*

We just successfully ran the engine's "Hello World." Here is the journey of that single log message:

1.  **The Event**: We shouted `Event{Type: "START_GAME"}`.
    *   *Analogy*: Someone pressed a button.
2.  **The Dispatcher**: The "Brain" heard it. It looked at its list of Rules (Behaviors).
3.  **The Trigger**: It checked `test_rule_1`.
    *   *Question*: "Does this rule care about 'START_GAME'?"
    *   *Answer*: **Yes.**
4.  **The Condition**: It checked logic.
    *   *Question*: "Is `AlwaysTrue` valid?"
    *   *Answer*: **Yes.**
5.  **The Action**: It executed the command.
    *   *Result*: `DebugLog` printed "Engine is running!"

### Why all this complexity for a print statement?
If we just wanted to print, we could have written `fmt.Println`.
But by building this **Pipeline**, we can now replace:
*   **Trigger**: "START_GAME" -> "PLAYER_HIT_ENEMY"
*   **Condition**: "AlwaysTrue" -> "EnemyHealth < 0"
*   **Action**: "Log" -> "PlaySound('boom.wav') AND SpawnParticles() AND AddScore(100)"

We created a generic machine that can handle *any* game logic without changing the hard code. We just feed it new Rules.

---

## 6. Architecture: The "Self-Contained" Object
*Why do we store the ID twice? (Once in the Map Key, once in the Struct)*

You noticed that `World.Entities` is a `map[string]*Entity`, but `Entity` also has an `ID` field.

**Isn't this redundant data?**
Yes.

**Why do we do it?**
1.  **Self-Knowledge**: If I pass an `*Entity` to a function (e.g., `CalculateDamage(e *Entity)`), that entity needs to know *who it is*.
    *   If the ID was only in the map key, the entity pointer `e` would be anonymous. We would have to pass `CalculateDamage(id string, e *Entity)` every time.
2.  **Reverse Lookup**: If you have the pointer, you can't find the map key without searching the entire map (which is slow O(N)).
3.  **Serialization**: When we save the game to JSON, we save the `[]Entity` list. The map is just a runtime lookup tool. The struct is the "source of truth."

**Rule of Thumb**:
It is okay to duplicate small data (like an ID string) if it makes the API significantly cleaner (avoiding extra arguments).

---

## 7. Interfaces vs Implementation: "The Contract"
*Understanding `behavior.Condition` vs `conditions.AlwaysTrue`*

We encountered a common Go pattern often confusing to newcomers coming from dynamic languages.

### The Question
"We have `behavior.go` defining a Condition, and `main.go` using a Condition. Why do they look different?"

### The Anatomy
1.  **The Interface (The "What")**
    *   *Location*: `behavior.go` -> `type Condition interface { Evaluate(...) bool }`
    *   *Purpose*: This is the **Contract**. It creates a standardized slot. It says "I don't care *what* object you give me, as long as it has an `Evaluate` method."
    *   *Analogy*: A generic power socket on the wall. It doesn't know if you'll plug in a toaster or a TV.

2.  **The Implementation (The "Which")**
    *   *Location*: `main.go` -> `conditions.AlwaysTrue{}`
    *   *Purpose*: This is the **Device**. It is a concrete struct that satisfies the contract.
    *   *Analogy*: The toaster. It has a plug that fits the socket.

### Why separate them?
By using the interface in `Behavior` (`Conditions []Condition`), we allow the engine to be **extensible**.
*   Today we use `AlwaysTrue`.
*   Tomorrow we write `IsEnemyVisible`.
*   The day after we write `HasEnoughMana`.

The generic `Behavior` code **never changes**. We just plug in different "devices" (Conditions).

---

## 8. The Dispatcher: The "God Object" Pattern
*Understanding `dispatcher := behavior.NewDispatcher(world, []*behavior.Behavior{testBehavior})`*

When you see this line in `main.go`, you're witnessing the creation of the engine's **Central Coordinator**.

### What is the Dispatcher?
```go
type Dispatcher struct {
    World      *entity.World      // The DATA (all entities)
    Behaviors  []*Behavior        // The RULES (all logic)
    eventQueue []core.Event       // The INBOX (pending events)
}
```

The Dispatcher is intentionally a **"God Object"** - it has access to everything. This seems to violate OOP principles, but it's a deliberate design choice.

### The Restaurant Manager Analogy
- **`World`** = The kitchen, tables, ingredients (the actual game state)
- **`Behaviors`** = The recipe book (rules for what to do)
- **`eventQueue`** = Customer orders coming in

The manager (Dispatcher) needs access to **all three** to coordinate the restaurant.

### Why This Breaks Traditional OOP (And Why That's OK)

**Traditional OOP teaches**: "Hide everything. Only expose what's necessary."

**Game engines need**: Everything is connected to everything.
- Enemy AI needs Player position
- UI needs Player health  
- Sound system needs Player state
- Physics needs Player velocity

If you follow strict encapsulation, you get **"getter hell"**:
```go
enemy.CheckDistance(player.GetPosition())
ui.Update(player.GetHealth())
sound.Play(player.GetWalkingState())
```

### The Controlled Global State Solution

The Dispatcher has access to everything, **but**:
1. **Single Responsibility**: Its only job is routing events (not rendering, physics, or input)
2. **Predictable Flow**: Events are processed one at a time, in order
3. **Controlled Interface**: Outside code can only `Emit()` events, not directly mutate the World

```
┌─────────────────────────────────┐
│         Dispatcher              │ ← The "Trusted Zone"
│  ┌──────┐      ┌──────────┐    │
│  │World │      │Behaviors │    │
│  └──────┘      └──────────┘    │
└─────────────────────────────────┘
         ↑
         │ Events only (controlled interface)
         │
    Game Loop / Input / Network
```

### The Constructor: Building the Brain

```go
dispatcher := behavior.NewDispatcher(world, []*behavior.Behavior{testBehavior})
```

This line:
1. **Creates** a new Dispatcher struct
2. **Gives it** the World (so behaviors can read/modify entities)
3. **Gives it** the Rulebook (the list of all active behaviors)

**Important**: This doesn't call the `Game` struct yet. That happens later:
```go
game := &Game{
    Dispatcher: dispatcher,  // ← Plugging the engine into the car
}
```

### Why One God Object is Better Than Many Small Objects

**Alternative (strict OOP)**:
```go
player.OnHit(enemy)
  → player.health -= damage
    → player.CheckDeath()
      → game.ShowGameOver()
        → ui.Display("You Died")
```
**Problems**: Deep call stacks, tight coupling, hard to test

**Event-Driven (our way)**:
```go
dispatcher.Emit(Event{Type: "PLAYER_HIT"})
  → HealthBehavior: Reduce health
  → DeathBehavior: Check if dead → Emit "PLAYER_DIED"
  → UIBehavior: Show death screen
  → SoundBehavior: Play death sound
```
**Benefits**: Flat structure, independent behaviors, easy to add/remove features

### The Golden Rule
Limit yourself to **one** God Object (the Dispatcher). Keep its logic simple (just routing). All game logic lives in **Behaviors**, which are small and testable.

This is why ECS (Entity-Component-System) and Event-Driven architectures dominate game development, even though they violate traditional OOP principles.

---

## 9. Payload: The Event's Cargo
*Understanding the data flow from `event.go` to `debug_log.go`*

When we ran the engine, we saw this output:
```
[ACTION] DebugLog: Engine is running! (Event: START_GAME)
        Payload: map[param:test_value]
```

Let's trace where this "Payload" comes from and what it means.

### What Does "Payload" Mean?

**Payload** is a term borrowed from shipping and aerospace:
- **Rocket**: The payload is the satellite (the useful cargo), not the fuel
- **Truck**: The payload is the boxes being delivered, not the truck itself
- **Letter**: The payload is the message inside, not the envelope

**In programming**: Payload = "The actual data being transmitted"

### The Definition: event.go

```go
// internal/engine/core/event.go
type Event struct {
    Type    EventType           // The "label" (what kind of event)
    A       string              // First entity involved (optional)
    B       string              // Second entity involved (optional)
    Payload map[string]any      // The "cargo" (flexible data)
}
```

**Key insight**: `Payload` is just a **field name** in the Event struct. Its type is `map[string]any`, which means:
- **Keys** are strings (e.g., `"damage"`, `"position"`, `"param"`)
- **Values** can be anything (`any` = any type: int, string, bool, struct, etc.)

This gives maximum flexibility!

### The Complete Data Flow

#### 1. **Creation** (main.go)
```go
dispatcher.Emit(core.Event{
    Type: "START_GAME",
    Payload: map[string]any{
        "param": "test_value",  // ← Data originates here
    },
})
```

#### 2. **Storage** (dispatcher.go)
```go
func (d *Dispatcher) Emit(ev core.Event) {
    d.eventQueue = append(d.eventQueue, ev)  // ← Queued for processing
}
```

#### 3. **Processing** (dispatcher.go)
```go
ctx := core.NewContext(d.World, ev)  // ← Event packed into Context
```

#### 4. **Delivery** (debug_log.go)
```go
func (a *DebugLog) Execute(ctx *core.Context) {
    // ctx.Event.Payload is now accessible!
    if ctx.Event.Payload != nil {
        fmt.Printf("\tPayload: %v\n", ctx.Event.Payload)  // ← Printed!
    }
}
```

### The Context: The Messenger Bag

The **Context** (`ctx`) is what carries data to actions:
```go
type Context struct {
    World *entity.World  // Access to all entities
    Event Event          // The event (including Payload)
}
```

Every action receives this Context, so it can:
- Read event data: `ctx.Event.Payload["damage"]`
- Access entities: `ctx.World.GetEntity(...)`
- Modify the game state

### Practical Examples

**Player takes damage:**
```go
Event{
    Type: "PLAYER_HIT",
    A: "player",
    B: "enemy_3",
    Payload: map[string]any{
        "damage": 25,
        "damageType": "fire",
    },
}
```

**Spawn an enemy:**
```go
Event{
    Type: "SPAWN_ENEMY",
    Payload: map[string]any{
        "x": 100,
        "y": 200,
        "enemyType": "goblin",
        "level": 5,
    },
}
```

**Key pressed:**
```go
Event{
    Type: "KEY_PRESSED",
    Payload: map[string]any{
        "key": "SPACE",
        "timestamp": 1234567890,
    },
}
```

### Why Separate Type and Payload?

- **Type**: Fast to check (string comparison) - used by Dispatcher to route events
- **Payload**: Flexible data container - used by Actions to get details

**Analogy**: Email
- **Subject line** = Type ("Meeting Reminder", "Invoice")
- **Email body** = Payload (the actual message content)

You scan subject lines quickly, but only read the body when relevant!

### The Takeaway

The Payload is the **flexible data container** that travels with every event. It's defined once in `event.go` but used everywhere in your engine to pass information between systems without tight coupling.
