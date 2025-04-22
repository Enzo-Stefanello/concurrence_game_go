package main

import (
	"math/rand"
	"time"
)

func GerarArmadilhas(jogo *Jogo, quantidade int, dano chan int) {
	for i := 0; i < quantidade; i++ {
		x := rand.Intn(len(jogo.Mapa[0]))
		y := rand.Intn(len(jogo.Mapa))

		if jogo.Mapa[y][x].tangivel {
			i--
			continue
		}

		invisivel := Elemento{
			simbolo:  ' ', // invisível
			cor:      CorVermelho,
			corFundo: CorPadrao,
			tangivel: false,
			visivel:  false,
		}
		jogo.Mapa[y][x] = invisivel

		go Armadilhas(jogo, x, y, dano) // Passa o canal de dano
	}
}

func Armadilhas(jogo *Jogo, x, y int, dano chan int) {
	ativa := false
	danoCausado := false // Flag para evitar causar dano repetido

	for {
		jogo.Mutex.Lock()
		distX := abs(jogo.PosX - x)
		distY := abs(jogo.PosY - y)

		// Quando o jogador se aproxima da armadilha (distância <= 3)
		if distX <= 3 && distY <= 3 {
			if !ativa {
				ativa = true
				jogo.Mapa[y][x].visivel = true
				jogo.Mapa[y][x].simbolo = 'X' // A armadilha se torna visível
			}
		} else {
			if ativa {
				ativa = false
				jogo.Mapa[y][x].visivel = false
				jogo.Mapa[y][x].simbolo = ' ' // A armadilha se torna invisível
			}
		}

		// Verifica se o jogador está pisando na armadilha
		if ativa && jogo.PosX == x && jogo.PosY == y {
			if !danoCausado {
				dano <- 10 // Envia 10 de dano
				danoCausado = true
			}
		} else {
			danoCausado = false // Reseta o flag quando o jogador se afasta da armadilha
		}

		jogo.Mutex.Unlock()
		time.Sleep(100 * time.Millisecond) // Faz a escuta contínua
	}
}
