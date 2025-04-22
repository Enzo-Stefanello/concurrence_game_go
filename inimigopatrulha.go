package main

import (
	"fmt"
	"math/rand"
	"time"
)

func comportamentoInimigo(jogo *Jogo, inimigo *Boss, done chan bool, perseguir chan bool, patrulhar chan bool) {
	for {
		select {
		case <-done:
			return
		case <-perseguir:
			jogo.Mutex.Lock()

			// Movimento em direção ao jogador
			var dx, dy int
			if jogo.PosX > inimigo.PosX {
				dx = 1
			} else if jogo.PosX < inimigo.PosX {
				dx = -1
			}
			if jogo.PosY > inimigo.PosY {
				dy = 1
			} else if jogo.PosY < inimigo.PosY {
				dy = -1
			}

			nx := inimigo.PosX + dx
			ny := inimigo.PosY + dy

			if nx == jogo.PosX && ny == jogo.PosY {
				jogo.VidaAtual -= 10
				jogo.StatusMsg = fmt.Sprintf("O inimigo causou 10 de dano! Vida atual: %d", jogo.VidaAtual)
			} else if jogoPodeMoverPara(jogo, nx, ny) {
				jogoMoverElemento(jogo, inimigo.PosX, inimigo.PosY, dx, dy)
				inimigo.PosX = nx
				inimigo.PosY = ny
			}

			interfaceDesenharJogo(jogo)
			jogo.Mutex.Unlock()

		case <-patrulhar:
			jogo.Mutex.Lock()

			// Movimento aleatório
			dx := rand.Intn(3) - 1
			dy := rand.Intn(3) - 1

			nx := inimigo.PosX + dx
			ny := inimigo.PosY + dy

			if nx == jogo.PosX && ny == jogo.PosY {
				jogo.VidaAtual -= 10
				jogo.StatusMsg = fmt.Sprintf("O inimigo causou 10 de dano! Vida atual: %d", jogo.VidaAtual)
			} else if jogoPodeMoverPara(jogo, nx, ny) {
				jogoMoverElemento(jogo, inimigo.PosX, inimigo.PosY, dx, dy)
				inimigo.PosX = nx
				inimigo.PosY = ny
			}

			interfaceDesenharJogo(jogo)
			jogo.Mutex.Unlock()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func monitorarJogador(jogo *Jogo, inimigo *Boss, perseguir chan bool, patrulhar chan bool, done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			jogo.Mutex.Lock()
			distX := abs(jogo.PosX - inimigo.PosX)
			distY := abs(jogo.PosY - inimigo.PosY)
			jogo.Mutex.Unlock()

			// Se o inimigo está próximo, envia para o canal 'perseguir', caso contrário, para 'patrulhar'
			if distX <= 3 && distY <= 3 {
				perseguir <- true
			} else {
				patrulhar <- true
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}
