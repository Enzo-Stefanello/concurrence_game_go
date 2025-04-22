package main

import (
	"fmt"
	//"math/rand"
	//"time"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1 // Move para cima
	case 'a':
		dx = -1 // Move para a esquerda
	case 's':
		dy = 1 // Move para baixo
	case 'd':
		dx = 1 // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, verifica se o jogador está interagindo com o portal
func personagemInteragir(jogo *Jogo) {
	if jogo.Mapa[jogo.PosY][jogo.PosX] == Armadilha {
		personagemTomarDano(jogo, 10) // Exemplo: dano de 10
	} else {
		jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
	}
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}

func personagemTomarDano(jogo *Jogo, dano int) {
	jogo.VidaAtual -= dano
	if jogo.VidaAtual < 0 {
		jogo.VidaAtual = 0
	}
	jogo.StatusMsg = fmt.Sprintf("Você tomou %d de dano! Vida restante: %d", dano, jogo.VidaAtual)
}
