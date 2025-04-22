package main

import "fmt"

// Função para gerenciar o dano recebido
func GerenciarDano(dano chan int, jogo *Jogo) {
	for {
		select {
		case d := <-dano:
			// Processa o dano recebido
			jogo.Mutex.Lock() // Protege o acesso à vida do jogador
			jogo.VidaAtual -= d
			if jogo.VidaAtual <= 0 {
				jogo.StatusMsg = "Game Over! Você morreu."
				jogo.Mutex.Unlock()
				interfaceDesenharJogo(jogo)
				return // Encerra a função ao atingir Game Over
			}
			jogo.StatusMsg = fmt.Sprintf("Você recebeu %d de dano! Vida restante: %d", d, jogo.VidaAtual)
			jogo.Mutex.Unlock()
			interfaceDesenharJogo(jogo) // Atualiza a tela após o dano
		}
	}
}
