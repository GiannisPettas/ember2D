package main

import (
	"image/color"
	"log"

	"github.com/GiannisPettas/ember2D/internal/engine/components"
	"github.com/GiannisPettas/ember2D/internal/engine/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Component managers (global for now)
var (
	world      *entity.World
	positions  *components.ComponentManager[components.Position]
	velocities *components.ComponentManager[components.Velocity]
	displays   *components.ComponentManager[components.Display]
)

type Game struct{}

func (g *Game) Update() error {
	// Movement system: move entities based on velocity
	velocities.Each(func(e entity.Entity, vel *components.Velocity) {
		if pos := positions.Get(e); pos != nil {
			pos.X += vel.X
			pos.Y += vel.Y
		}
	})

	// Bounce enemies off screen edges
	for _, e := range world.Tags().GetEntitiesByTag("enemy") {
		if pos := positions.Get(e); pos != nil {
			if vel := velocities.Get(e); vel != nil {
				if pos.Y > 440 || pos.Y < 0 {
					vel.Y = -vel.Y
				}
				if pos.X > 600 || pos.X < 0 {
					vel.X = -vel.X
				}
			}
		}
	}

	// Player movement with arrow keys
	for _, e := range world.Tags().GetEntitiesByTag("player") {
		if pos := positions.Get(e); pos != nil {
			speed := 3.0
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
				pos.Y -= speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
				pos.Y += speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
				pos.X -= speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
				pos.X += speed
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Render system: draw all entities with Position + Display
	displays.Each(func(e entity.Entity, d *components.Display) {
		if pos := positions.Get(e); pos != nil {
			clr := color.RGBA{R: d.R, G: d.G, B: d.B, A: 255}
			ebitenutil.DrawRect(screen, pos.X, pos.Y, d.Width, d.Height, clr)
		}
	})

	// Debug info
	ebitenutil.DebugPrint(screen, "ember2D - Arrow keys to move")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	// Initialize world and component managers
	world = entity.NewWorld()
	positions = components.NewComponentManager[components.Position]()
	velocities = components.NewComponentManager[components.Velocity]()
	displays = components.NewComponentManager[components.Display]()

	// Create player (blue square)
	player := world.CreateEntity("player")
	positions.Add(player, components.Position{X: 300, Y: 220})
	displays.Add(player, components.Display{Width: 30, Height: 30, R: 50, G: 100, B: 255})

	// Create enemies (red squares that bounce around)
	for i := 0; i < 5; i++ {
		enemy := world.CreateEntity("enemy")
		positions.Add(enemy, components.Position{
			X: float64(80 + i*100),
			Y: float64(50 + i*60),
		})
		velocities.Add(enemy, components.Velocity{
			X: float64(1 + i),
			Y: float64(2 - i),
		})
		displays.Add(enemy, components.Display{Width: 20, Height: 20, R: 255, G: 50, B: 50})
	}

	// Run game
	ebiten.SetWindowTitle("ember2D Runtime")
	ebiten.SetWindowSize(640, 480)
	log.Println("Starting ember2D runtime...")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
