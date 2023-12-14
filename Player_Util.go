package main

import (
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type player struct {
	spriteSheet     *ebiten.Image
	xLoc, yLoc      int
	playerState     int
	direction       int
	fr, frD         int
	SPD, ATK, SCORE int
	sword           bool
	notif           int
	health          float64
}

func drawPlayer(screen *ebiten.Image, p player) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.xLoc), float64(p.yLoc))
	screen.DrawImage(p.spriteSheet.SubImage(image.Rect(p.fr*player_width,
		p.direction*player_height, p.fr*player_width+player_width,
		p.direction*player_height+player_height)).(*ebiten.Image), op)
}

func getPlayerInput(game *mapGame) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && game.player.xLoc > 0 {
		game.player.direction = LEFT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) &&
		game.player.xLoc < 960-player_width {
		game.player.direction = RIGHT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) &&
		game.player.yLoc > 0 {
		game.player.direction = UP
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
		game.player.yLoc < 960-player_height {
		game.player.direction = DOWN
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) && game.player.sword == true {
		game.player.direction = SPACE
	}
}

func handlePlayerMovement(player *player, game *tileMap) {
	//fmt.Println("Xloc:", player.xLoc, "Yloc:", player.yLoc)
	switch player.direction {
	case LEFT:
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) &&
			(!game.isTileBarrier(player.xLoc+20, player.yLoc-20) || !game.isTileBarrier(player.xLoc+20, player.yLoc+20) ||
				ebiten.IsKeyPressed(ebiten.KeyJ)) {
			player.xLoc -= player.SPD
			if player.xLoc <= 0 {
				player.xLoc += player.SPD
			}
		}
	case RIGHT:
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) &&
			(!game.isTileBarrier(player.xLoc+player_width-20, player.yLoc+20) || ebiten.IsKeyPressed(ebiten.KeyJ)) {
			player.xLoc += player.SPD
			if player.xLoc >= 960 {
				player.xLoc -= player.SPD
			}
		}
	case DOWN:
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) &&
			(!game.isTileBarrier(player.xLoc+40, player.yLoc+player_height) || ebiten.IsKeyPressed(ebiten.KeyJ)) {
			player.yLoc += player.SPD
			if player.yLoc >= 900 {
				player.yLoc -= player.SPD
			}
		}
	case UP:
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) &&
			(!game.isTileBarrier(player.xLoc+40, player.yLoc+player_height-30) || ebiten.IsKeyPressed(ebiten.KeyJ)) {
			player.yLoc -= player.SPD
			if player.yLoc <= 0 {
				player.yLoc += player.SPD
			}
		}
	}
}

func checkCollisions(player player, npc npc, height, width int) bool {
	playerBounds := collision.BoundingBox{
		X:      float64(player.xLoc),
		Y:      float64(player.yLoc),
		Width:  float64(player_width),
		Height: float64(player_height),
	}
	npcBounds := collision.BoundingBox{
		X:      float64(npc.xLoc),
		Y:      float64(npc.yLoc),
		Width:  float64(width),
		Height: float64(height),
	}
	if collision.AABBCollision(playerBounds, npcBounds) {
		return true
	}
	return false
}
