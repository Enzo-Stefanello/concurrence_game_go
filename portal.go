package main

import (
	"math/rand"
	"time"
)

func gerarPosicaoAleatoriaVazia(jogo *Jogo) (int, int) {
	//rand.Seed(time.Now().UnixNano())
	for {
		x := rand.Intn(len(jogo.Mapa[0]))
		y := rand.Intn(len(jogo.Mapa))
		if jogo.Mapa[y][x] == Vazio {
			return x, y
		}
	}
}

// Função para gerar posições aleatórias para o portal
/*func gerarPosicaoPortal(jogo *Jogo) (int, int) {
	// Gera posições aleatórias dentro dos limites do mapa
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len(jogo.Mapa[0])), rand.Intn(len(jogo.Mapa))
}*/

// Função para gerar portais aleatórios sem fechar o canal diretamente após interação
func GerarPortaisAleatorios(jogo *Jogo, numPortais int, duracao time.Duration, done chan bool) {
	for i := 0; i < numPortais; i++ {
		go func() {
			// Gera uma posição aleatória para o portal
			x, y := gerarPosicaoAleatoriaVazia(jogo)

			// Coloca o portal no mapa
			jogo.Mutex.Lock()
			jogo.Mapa[y][x] = Portal
			jogo.Mutex.Unlock()

			// Espera o tempo do portal durar
			tempo := time.After(duracao)
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-tempo:
					// Se o tempo expirar, remove o portal
					jogo.Mutex.Lock()
					if jogo.Mapa[y][x] == Portal {
						jogo.Mapa[y][x] = Vazio
					}
					jogo.Mutex.Unlock()
					return
				case <-ticker.C:
					// Verifica se o jogador está no portal
					jogo.Mutex.Lock()
					if jogo.PosX == x && jogo.PosY == y {
						// Remove o portal
						jogo.Mapa[y][x] = Vazio
						// Teleporta para nova posição válida aleatória
						nx, ny := gerarPosicaoAleatoriaVazia(jogo)
						jogo.PosX = nx
						jogo.PosY = ny
						jogo.StatusMsg = "Você foi teletransportado para uma posição aleatória!"

						// Apenas atualiza o estado do jogo, sem fechar o canal
						jogo.Mutex.Unlock()
						return
					}
					jogo.Mutex.Unlock()
				}
			}
		}()
	}
}

// Função principal ou outra parte do código onde o canal é fechado
func FinalizarJogo(done chan bool) {
	// Após todas as goroutines estarem completas, você pode fechar o canal
	close(done)
}
