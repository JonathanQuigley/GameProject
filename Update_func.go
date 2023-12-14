package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

func (game *mapGame) Update() error {
	if game.player.health <= 0 {
		game.player.playerState = 0
		game.player.health = float64(100 + game.player.SCORE/2)
		game.player.SPD = int(2 + float64(game.player.SCORE/100))
		game.player.SCORE = 0
		game.player.xLoc = 580
		game.player.yLoc = 500
	}
	getPlayerInput(game)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		game.player.frD += 1
		if game.player.frD%player_fps == 0 {
			game.player.fr += 1
			if game.player.fr >= player_frame_per_sheet {
				game.player.fr = 0
				game.sounds.playerwalk.Rewind()
				game.sounds.playerwalk.Play()
			}
		}
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) && game.player.sword == false {
		game.sounds.playerNoSwordATK.Rewind()
		game.sounds.playerNoSwordATK.Play()
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) && game.player.sword == true {
		for i := 0; i < 5; i++ {
			game.player.frD += 1
			if game.player.frD%20 == 0 {
				game.player.fr += 1
				if game.player.fr >= 4 {
					game.player.direction = 0
					game.player.fr = 0
					game.sounds.playerSwordATK.Rewind()
					game.sounds.playerSwordATK.Play()
				}
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.Key1) {
		game.player.notif = 0
	}
	if game.player.playerState == 0 {
		//START MAP

		updateNPC(&game.gmWar, 8, 9, 0, game.gmWar.xLoc, game.gmWar.yLoc, 0, 0)
		updateNPC(&game.mimic, 16, 4, 0, game.mimic.xLoc, game.mimic.yLoc, 0, 0)
		updateNPC(&game.swords, 16, 6, 5, game.swords.xLoc, game.swords.yLoc, 0, 0)
		handlePlayerMovement(&game.player, &game.map1)
		if checkCollisions(game.player, game.swords, sword_width, sword_height) {
			if ebiten.IsKeyPressed(ebiten.KeySpace) && game.player.sword == false {
				game.swords.xLoc = -1000
				game.swords.yLoc = -1000
				game.player.ATK *= 3
				game.player.notif = 1
				game.player.sword = true
				game.sounds.NPCtalk.Rewind()
				game.sounds.NPCtalk.Play()
			}
		} else if checkCollisions(game.player, game.mimic, mimic_width, mimic_height) {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				game.mimic.xLoc = -1000
				game.mimic.yLoc = -1000
				game.player.SPD += 2
				game.player.notif = 2
				game.sounds.mimicHit.Rewind()
				game.sounds.mimicHit.Play()
				game.sounds.mimicDead.Rewind()
				game.sounds.mimicDead.Play()
			}
		} else if checkCollisions(game.player, game.gmWar, mimic_width, mimic_height) {
			game.player.notif = 3
		} else if ebiten.IsKeyPressed(ebiten.KeyC) {
			game.player.notif = 4
			game.sounds.NPCtalk.Rewind()
			game.sounds.NPCtalk.Play()
		} else if ebiten.IsKeyPressed(ebiten.KeyX) {
			game.player.notif = 5
		}
		if game.player.xLoc >= 880 && (game.player.yLoc >= 240 && game.player.yLoc <= 360) {
			game.sounds.playerExit.Rewind()
			game.sounds.playerExit.Play()
			game.player.playerState = 1
			game.player.xLoc = 80
			game.player.yLoc = 500
		}
	} else if game.player.playerState == 1 {
		//INTERMEDIARY MAP

		updateNPC(&game.mystic, 16, 0, 4, game.mystic.xLoc, game.mystic.yLoc, 0, 0)
		updateNPC(&game.demon, 8, 6, 0, game.demon.xLoc, game.demon.yLoc, 0, 0)
		handlePlayerMovement(&game.player, &game.map2)
		if game.player.xLoc <= 30 && (game.player.yLoc >= 440 && game.player.yLoc <= 550) {
			game.sounds.playerExit.Rewind()
			game.sounds.playerExit.Play()
			game.player.playerState = 0
			game.player.xLoc = 860
			game.player.yLoc = 340
		} else if game.player.yLoc <= 10 && (game.player.xLoc >= 390 && game.player.xLoc <= 470) {
			game.sounds.playerExit.Rewind()
			game.sounds.playerExit.Play()
			game.player.playerState = 2
			game.player.xLoc = 380
			game.player.yLoc = 840
		}
	} else if game.player.playerState == 2 {
		//BATTLE MAP

		updateNPC(&game.enemyskele, 8, 4, 0, rand.Intn(500), 100, 450, 3)
		handlePanicSkele(game)
		handleDragonState(game)
		handlePlayerMovement(&game.player, &game.map3)
		if checkCollisions(game.player, game.enemyskele, 48, 48) && game.player.sword == true {
			game.player.health -= 0.1
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				game.sounds.enemyHit.Rewind()
				game.sounds.enemyHit.Play()
				game.player.SCORE += 1
				game.enemyskele.yLoc = 100
				game.enemyskele.xLoc = rand.Intn(500)
			}
		}
		if checkCollisions(game.player, game.dragon, 48, 48) {
			game.player.health -= 0.5
		}
		if checkCollisions(game.player, game.dragon, 80, 80) && game.player.sword == true {
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				game.sounds.enemyHit.Rewind()
				game.sounds.enemyHit.Play()
				game.dragon.health -= float64(game.player.ATK)
				if game.dragon.health <= 0 {
					game.dragon.xLoc = 10
					game.dragon.yLoc = 100
					game.dragon.state = 0
					game.player.SCORE += 100
				}
			}
		}
		if game.player.yLoc >= 860 && (game.player.xLoc >= 310 && game.player.xLoc <= 410) {
			game.sounds.playerExit.Rewind()
			game.sounds.playerExit.Play()
			game.player.playerState = 1
			game.player.xLoc = 440
			game.player.yLoc = 30
		}

	}
	return nil
}
