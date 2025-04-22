// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"sync"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
	visivel  bool // Variável para controlar a visibilidade
}

// Boss representa um inimigo no jogo
type Boss struct {
	PosX int // Coordenada X do inimigo
	PosY int // Coordenada Y do inimigo
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	Vida           int
	VidaMaxima     int
	PosX, PosY     int        // posição atual do personagem principal
	Inimigos       []Boss     // lista de inimigos no jogo
	UltimoVisitado Elemento   // elemento que estava na posição do personagem antes de mover
	StatusMsg      string     // mensagem para a barra de status
	Mutex          sync.Mutex // Protege o acesso ao mapa e ao estado do jogo
	VidaAtual      int
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false, true}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false, true}
	Portal     = Elemento{'O', CorAzul, CorPadrao, false, true} // Novo elemento para o portal
	Armadilha  = Elemento{'X', CorVermelho, CorPadrao, true, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// Inicializa o jogo com o último elemento visitado como vazio e uma lista vazia de inimigos
	return Jogo{
		UltimoVisitado: Vazio,
		Inimigos:       []Boss{}, // Inicializa a lista de inimigos vazia
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
				jogo.Inimigos = append(jogo.Inimigos, Boss{PosX: x, PosY: y})
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y
			case Portal.simbolo: // Novo caso para o portal
				e = Portal
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
// Verifica se o jogador pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem o inimigo na posição atual
	elemento := jogo.Mapa[y][x]

	// Se o elemento for um inimigo, garante que ele mantém suas propriedades
	if elemento.simbolo == Inimigo.simbolo {
		jogo.Mapa[y][x] = jogo.UltimoVisitado
		jogo.UltimoVisitado = jogo.Mapa[ny][nx]
		jogo.Mapa[ny][nx] = Inimigo // Garante que a posição de destino é sempre um inimigo
	}
}

// filepath: c:\Users\enzos\Desktop\fppd\jogo.go
func inicializarInimigos(jogo *Jogo, numInimigos int) {
	for i := 0; i < numInimigos; i++ {
		x, y := gerarPosicaoAleatoriaVazia(jogo)
		inimigo := Boss{PosX: x, PosY: y}
		jogo.Inimigos = append(jogo.Inimigos, inimigo)
		jogo.Mapa[y][x] = Inimigo // Representa o inimigo no mapa
	}
}
