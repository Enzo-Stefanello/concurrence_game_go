package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorVerde           = termbox.ColorGreen
	CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
	CorAzul            = termbox.ColorBlue
	CorAmarelo         = termbox.ColorYellow
	CorCiano           = termbox.ColorCyan
	CorMagenta         = termbox.ColorMagenta
	CorBranco          = termbox.ColorWhite
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	if jogo.VidaAtual <= 0 {
		interfaceExibirGameOver()
		interfaceFinalizar()
		panic("Game Over!") // Encerra o jogo
	}

	interfaceLimparTela()

	// Renderiza apenas o valor numérico da vida
	renderizarBarraDeVida(jogo.VidaAtual, jogo.VidaMaxima)

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y+2, elem) // Ajusta o mapa para começar abaixo da vida
		}
	}

	// Desenha o personagem sobre o mapa
	interfaceDesenharElemento(jogo.PosX, jogo.PosY+2, Personagem) // Ajusta a posição do personagem

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo)

	// Força a atualização do terminal
	interfaceAtualizarTela()
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	linhaStatus := 2 + len(jogo.Mapa) + 1 // Linha abaixo do mapa
	linhaInstrucoes := linhaStatus + 1    // Instruções ainda mais abaixo

	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, linhaStatus, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover e E para interagir. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, linhaInstrucoes, c, CorTexto, CorPadrao)
	}
}

func renderizarBarraDeVida(vidaAtual, vidaMaxima int) {
	// Exibe o texto "Vida: X/Y" na linha 0
	texto := fmt.Sprintf("Vida: %d/%d", vidaAtual, vidaMaxima)
	for i, c := range texto {
		termbox.SetCell(i, 0, c, CorTexto, CorPadrao) // Linha 0 para o texto
	}
}

func interfaceExibirGameOver() {
	interfaceLimparTela()

	mensagem := "GAME OVER"
	instrucao := "Pressione ESC para sair."

	// Centraliza a mensagem na tela
	largura, altura := termbox.Size()
	x := (largura - len(mensagem)) / 2
	y := altura / 2

	for i, c := range mensagem {
		termbox.SetCell(x+i, y, c, CorVermelho, CorPadrao)
	}

	for i, c := range instrucao {
		termbox.SetCell(x+i, y+1, c, CorTexto, CorPadrao)
	}

	interfaceAtualizarTela()

	// Aguarda o jogador pressionar ESC para sair
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break
		}
	}
}
