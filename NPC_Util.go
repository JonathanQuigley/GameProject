package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
)

type npc struct {
	spritesheet *ebiten.Image
	xLoc, yLoc  int
	direction   int
	fr, frD     int
	health      float64
	state       int
}

//STATE -> 0 = patrol or just moving | 1 = Go after player | 2 = Death, loop back to state 0

func drawNPC(screen *ebiten.Image, n npc, width, height int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.xLoc), float64(n.yLoc))
	screen.DrawImage(n.spritesheet.SubImage(image.Rect(n.fr*width,
		n.direction*height, n.fr*width+width,
		n.direction*height+height)).(*ebiten.Image), op)
}

func updateNPC(n *npc, fps, frames, dir, x, y, y1, sp int) {
	n.frD++
	if n.frD%fps == 0 {
		n.yLoc += sp
		n.fr++
		if n.fr >= frames {
			n.fr = 0
			n.direction++
			if n.direction >= dir {
				n.direction = 0
			}
		} else if n.yLoc >= y1 {
			n.yLoc = y
			n.xLoc = x
		}
	}
}

func initNPC(imgName string, x, y, z int) npc {
	npcImg := LoadEmbeddedImage("NPCs", imgName)
	return npc{
		spritesheet: npcImg,
		xLoc:        x,
		yLoc:        y,
		state:       0,
		health:      float64(z),
	}
}

func handlePanicSkele(game *mapGame) {
	game.panicSkele.frD += 1
	if game.panicSkele.frD%8 == 0 {
		switch game.panicSkele.direction {
		case 0:
			game.panicSkele.yLoc -= 4
			game.panicSkele.fr += 1
			if game.panicSkele.yLoc <= 600 {
				game.panicSkele.direction = 1
			} else if game.panicSkele.fr >= 3 {
				game.panicSkele.fr = 0
			}
		case 1:
			game.panicSkele.xLoc += 4
			game.panicSkele.fr += 1
			if game.panicSkele.xLoc >= 200 {
				game.panicSkele.direction = 2
			} else if game.panicSkele.fr >= 3 {
				game.panicSkele.fr = 0
			}
		case 2:
			game.panicSkele.yLoc += 4
			game.panicSkele.fr += 1
			if game.panicSkele.yLoc >= 800 {
				game.panicSkele.direction = 3
			} else if game.panicSkele.fr >= 3 {
				game.panicSkele.fr = 0
			}
		case 3:
			game.panicSkele.xLoc -= 4
			game.panicSkele.fr += 1
			if game.panicSkele.xLoc <= 10 {
				game.panicSkele.direction = 0
			} else if game.panicSkele.fr >= 3 {
				game.panicSkele.fr = 0
			}
		}
	}
}

func handleDragonPatrol(game *mapGame) {
	game.dragon.frD += 1
	if game.dragon.frD%8 == 0 {
		switch game.dragon.direction {
		case 0:
			game.dragon.yLoc -= 6
			game.dragon.fr += 1
			if game.dragon.yLoc <= 10 {
				game.dragon.direction = 1
			} else if game.dragon.fr >= 3 {
				game.dragon.fr = 0
			}
		case 1:
			game.dragon.xLoc += 6
			game.dragon.fr += 1
			if game.dragon.xLoc >= 820 {
				game.dragon.direction = 2
			} else if game.dragon.fr >= 3 {
				game.dragon.fr = 0
			}
		case 2:
			game.dragon.yLoc += 6
			game.dragon.fr += 1
			if game.dragon.yLoc >= 420 {
				game.dragon.direction = 3
			} else if game.dragon.fr >= 3 {
				game.dragon.fr = 0
			}
		case 3:
			game.dragon.xLoc -= 6
			game.dragon.fr += 1
			if game.dragon.xLoc <= 10 {
				game.dragon.direction = 0
			} else if game.dragon.fr >= 3 {
				game.dragon.fr = 0
			}
		}
	}
}

func handleDragonState(game *mapGame) {
	switch game.dragon.state {
	case 0:
		handleDragonPatrol(game)
		if game.player.yLoc <= 420 {
			game.dragon.state = 1
		}
	case 1:
		startRow := game.dragon.yLoc / game.map3.Level.TileHeight
		startCol := game.dragon.xLoc / game.map3.Level.TileWidth
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			startCell := game.map3.pathMap.Get(startCol, startRow)
			endcol := game.player.xLoc / game.map3.Level.TileWidth
			endrow := game.player.yLoc / game.map3.Level.TileHeight
			endCell := game.map3.pathMap.Get(endcol, endrow)
			game.map3.path = game.map3.pathMap.GetPathFromCells(startCell, endCell, false, false)
		}
		if game.map3.path != nil {
			pathCell := game.map3.path.Current()
			if math.Abs(float64(pathCell.X*game.map3.Level.TileWidth)-float64(game.dragon.xLoc)) <= 2 &&
				math.Abs(float64(pathCell.Y*game.map3.Level.TileHeight)-float64(game.dragon.yLoc)) <= 2 { //if we are now on the tile we need to be on
				game.map3.path.Advance()
			}
			direction := 0.0
			if pathCell.X*game.map3.Level.TileWidth > game.dragon.xLoc {
				direction = 1.0
			} else if pathCell.X*game.map3.Level.TileWidth < game.dragon.xLoc {
				direction = -1.0
			}
			Ydirection := 0.0
			if pathCell.Y*game.map3.Level.TileHeight > game.dragon.yLoc {
				Ydirection = 1.0
			} else if pathCell.Y*game.map3.Level.TileHeight < game.dragon.yLoc {
				Ydirection = -1.0
			}
			if int(Ydirection) < 0 {
				game.dragon.yLoc -= 1
				game.dragon.frD += 1
				game.dragon.direction = 0
				if game.dragon.frD%16 == 0 {
					game.dragon.fr += 1
					if game.dragon.fr >= 3 {
						game.dragon.fr = 0
					}
				}
			} else if int(Ydirection) > 0 {
				game.dragon.yLoc += 1
				game.dragon.frD += 1
				game.dragon.direction = 2
				if game.dragon.frD%16 == 0 {
					game.dragon.fr += 1
					if game.dragon.fr >= 3 {
						game.dragon.fr = 0
					}
				}
			}
			if int(direction) < 0 {
				game.dragon.xLoc -= 1
				game.dragon.frD += 1
				game.dragon.direction = 3
				if game.dragon.frD%16 == 0 {
					game.dragon.fr += 1
					if game.dragon.fr >= 3 {
						game.dragon.fr = 0
					}
				}
			} else if int(direction) > 0 {
				game.dragon.xLoc += 1
				game.dragon.frD += 1
				game.dragon.direction = 1
				if game.dragon.frD%16 == 0 {
					game.dragon.fr += 1
					if game.dragon.fr >= 3 {
						game.dragon.fr = 0
					}
				}
			} else if checkCollisions(game.player, game.dragon, 48, 48) {
				game.dragon.frD += 1
				game.dragon.direction = 2
				if game.dragon.frD%16 == 0 {
					game.dragon.fr += 1
					if game.dragon.fr >= 3 {
						game.dragon.fr = 0
					}
				}
			}
			if game.player.yLoc >= 440 || game.dragon.yLoc >= 440 {
				game.dragon.state = 2
			}
		}
	case 2:
		game.dragon.yLoc -= 6
		if game.dragon.yLoc <= 20 {
			game.dragon.state = 0
		}
	}
}
