package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func (game mapGame) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	switch game.player.playerState {
	case 0:
		drawTileMap(screen, game.map1, op)
		drawPlayer(screen, game.player)
		drawNPC(screen, game.gmWar, player_width, player_height)
		drawNPC(screen, game.mimic, mimic_width, mimic_height)
		drawNPC(screen, game.swords, sword_width, sword_height)
		drawTexts(screen, game.font, game.player.notif)
		DrawCenteredText(screen, game.font, "START MAP", 480, 40)
		DrawCenteredText(screen, game.font2, controls, 820, 50)
		DrawCenteredText(screen, game.font2, objectives, 810, 100)
		DrawCenteredText(screen, game.font2, notes, 790, 170)
		scoreStr := fmt.Sprintf("Score: %d\nHealth: %f", game.player.SCORE, game.player.health)
		ebitenutil.DebugPrint(screen, scoreStr)
	case 1:
		drawTileMap(screen, game.map2, op)
		drawPlayer(screen, game.player)
		drawNPC(screen, game.mystic, player_width, player_height)
		drawNPC(screen, game.demon, demon_id_wt, demon_id_ht)
		drawTexts(screen, game.font, 0)
		DrawCenteredText(screen, game.font, "INTERMEDIARY MAP", 200, 40)
		scoreStr := fmt.Sprintf("Score: %d\nHealth: %f", game.player.SCORE, game.player.health)
		ebitenutil.DebugPrint(screen, scoreStr)
	case 2:
		drawTileMap(screen, game.map3, op)
		drawPlayer(screen, game.player)
		drawNPC(screen, game.panicSkele, skele_walk_wt, skele_walk_ht)
		drawNPC(screen, game.enemyskele, player_width, player_height)
		drawNPC(screen, game.dragon, dragon_width, dragon_height)
		drawTexts(screen, game.font, game.player.notif)
		DrawCenteredText(screen, game.font, "BATTLE MAP", 480, 40)
		scoreStr := fmt.Sprintf("Score: %d\nHealth: %f", game.player.SCORE, game.player.health)
		ebitenutil.DebugPrint(screen, scoreStr)
	}
}

func (game mapGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// ACCESSORY FUNCTIONS
func drawTileMap(screen *ebiten.Image, tileMap tileMap, op *ebiten.DrawImageOptions) {
	for tileY := 0; tileY < tileMap.Level.Height; tileY += 1 {
		for tileX := 0; tileX < tileMap.Level.Width; tileX += 1 {
			op.GeoM.Reset()
			TileXpos := float64(tileMap.Level.TileWidth * tileX)
			TileYpos := float64(tileMap.Level.TileHeight * tileY)
			op.GeoM.Translate(TileXpos, TileYpos)
			tileToDraw := tileMap.Level.Layers[0].Tiles[tileY*tileMap.Level.Width+tileX]
			ebitenTileToDraw := tileMap.tileHash[tileToDraw.ID]
			screen.DrawImage(ebitenTileToDraw, op)
		}
	}
}

func drawTexts(screen *ebiten.Image, f font.Face, notif int) {
	texts := []string{
		"   ",
		"SWORD ACQUIRED ATK +20 (press 1 to hide)",
		"MIMIC DEFEATED SPD DOUBLED (press 1 to hide)",
		"Press C to talk to Warlock",
		"You must continue your journey... (press X)",
		"Press onward and Slay your enemies! (Press 1)",
	}

	if notif >= 1 && notif < len(texts) {
		DrawCenteredText(screen, f, texts[notif], 480, 380)
	}
}

func LoadScoreFont() font.Face {
	//originally inspired by https://www.fatoldyeti.com/posts/roguelike16/
	trueTypeFont, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		fmt.Println("Error loading font for score:", err)
	}
	fontFace, err := opentype.NewFace(trueTypeFont, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println("Error loading font of correct size for score:", err)
	}
	return fontFace
}

func LoadControlsfont() font.Face {
	//originally inspired by https://www.fatoldyeti.com/posts/roguelike16/
	trueTypeFont, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		fmt.Println("Error loading font for score:", err)
	}
	fontFace, err := opentype.NewFace(trueTypeFont, &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println("Error loading font of correct size for score:", err)
	}
	return fontFace
}

func DrawCenteredText(screen *ebiten.Image, font font.Face, s string, cx, cy int) { //from https://github.com/sedyh/ebitengine-cheatsheet
	bounds := text.BoundString(font, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, colornames.White)
}
