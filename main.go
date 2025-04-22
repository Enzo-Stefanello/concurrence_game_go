package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	rand.Seed(time.Now().UnixNano())
	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := inicializarJogo()

	// Carrega o mapa
	if err := jogoCarregarMapa(mapaFile, jogo); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar o mapa: %v\n", err)
		os.Exit(1)
	}

	// Adiciona inimigos ao jogo
	inicializarInimigos(jogo, 5) // Adiciona 5 inimigos ao jogo

	// Canais para comunicação
	dano := make(chan int)
	done := make(chan bool)
	//perseguir := make(chan bool)
	//patrulhar := make(chan bool)

	// Inicia elementos concorrentes
	wg.Add(1)
	go func() {
		defer wg.Done()
		GerarPortaisAleatorios(jogo, 3, 30*time.Second, done) // 3 portais, durando 30 segundos
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GerarArmadilhas(jogo, 5, dano) // Agora passa o canal de dano
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GerenciarDano(dano, jogo)
	}()

	// Inicia o comportamento de inimigos (movimento aleatório ou perseguição)
	// Inicia o comportamento de inimigos com canais dedicados
	for i := range jogo.Inimigos {
		inimigo := &jogo.Inimigos[i]

		perseguir := make(chan bool)
		patrulhar := make(chan bool)
		// Goroutine do comportamento do inimigo
		wg.Add(1)
		go func(inimigo *Boss, perseguir, patrulhar chan bool) {
			defer wg.Done()
			comportamentoInimigo(jogo, inimigo, done, perseguir, patrulhar)
		}(inimigo, perseguir, patrulhar)

		// Goroutine do monitor que decide o comportamento
		wg.Add(1)
		go func(inimigo *Boss, perseguir, patrulhar chan bool) {
			defer wg.Done()
			monitorarJogador(jogo, inimigo, perseguir, patrulhar, done)
		}(inimigo, perseguir, patrulhar)
	}

	// Loop principal
loop:
	for {
		// Processa entrada do jogador
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, jogo); !continuar {
			break loop
		}

		// Redesenha o jogo
		interfaceDesenharJogo(jogo)

		// Aguarda um curto intervalo antes de continuar
		time.Sleep(100 * time.Millisecond)
	}

	// Finaliza goroutines
	close(done)
	wg.Wait()

	// Limpa o terminal antes de exibir a mensagem final
	interfaceLimparTela()
	interfaceAtualizarTela()

	// Exibe as mensagens de encerramento
	fmt.Println("Encerrando o jogo...")
	fmt.Println("Jogo encerrado com sucesso.")
}

// Inicializa o jogo com valores padrão
func inicializarJogo() *Jogo {
	jogo := jogoNovo()       // Supondo que jogoNovo() retorna Jogo, não *Jogo
	jogoPtr := &jogo         // Cria um ponteiro para o valor retornado
	jogoPtr.VidaAtual = 100  // Inicializa a vida atual
	jogoPtr.VidaMaxima = 100 // Inicializa a vida máxima
	jogoPtr.StatusMsg = "Bem-vindo ao jogo"
	return jogoPtr
}
