package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width      = 20  // Largura da tela do jogo
	height     = 10  // Altura da tela do jogo
	borderChar = '▓' // Caractere usado para desenhar as bordas da tela
)

// Definição do tipo 'point' para representar coordenadas (x, y) na tela
type point struct {
	x, y int
}

// Constantes para representar direções
const (
	up    = iota // 0
	down         // 1
	left         // 2
	right        // 3
)

var (
	snake       []point                  // Slice para armazenar as coordenadas da cobra
	food        point                    // Coordenadas da comida
	direction   = right                  // Direção inicial da cobra
	gameOver    = false                  // Flag para indicar o fim do jogo
	score       = 0                      // Pontuação do jogador
	refreshRate = 100 * time.Millisecond // Taxa de atualização do jogo
)

// Função principal
func main() {
	// Inicializa a biblioteca termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc) // Configura o modo de entrada para aceitar teclas de escape

	snake = []point{{width / 2, height / 2}} // Inicializa a cobra no centro da tela
	placeFood()                              // Coloca a comida em uma posição aleatória

	go gameLoop() // Inicia o loop principal do jogo em uma goroutine

	// Loop de evento principal
	for {
		ev := termbox.PollEvent() // Aguarda eventos de entrada do usuário
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if direction != down {
					direction = up
				}
			case termbox.KeyArrowDown:
				if direction != up {
					direction = down
				}
			case termbox.KeyArrowLeft:
				if direction != right {
					direction = left
				}
			case termbox.KeyArrowRight:
				if direction != left {
					direction = right
				}
			case termbox.KeyEsc:
				gameOver = true // Define a flag de fim de jogo como verdadeira
				return
			case termbox.KeyEnter:
				if gameOver {
					// Reinicia o jogo se estiver no estado de jogo encerrado
					snake = []point{{width / 2, height / 2}}
					placeFood()
					direction = right
					gameOver = false
					score = 0
					go gameLoop() // Inicia um novo loop de jogo em uma goroutine
				}
			}
		}
	}
}

// Função gameLoop controla a lógica do jogo e a taxa de atualização
func gameLoop() {
	ticker := time.NewTicker(refreshRate) // Cria um ticker para controlar a taxa de atualização
	defer ticker.Stop()                   // Garante que o ticker seja parado quando a função retornar

	// Loop do jogo
	for !gameOver {
		select {
		case <-ticker.C: // Aguarda o próximo tick do ticker
			if !gameOver {
				update() // Atualiza a lógica do jogo
				render() // Renderiza a tela do jogo
			}
		}
	}
}

// Função update atualiza o estado da cobra e verifica as colisões
func update() {
	if gameOver {
		return // Se o jogo acabou, não faça nada
	}

	head := snake[0] // Obtém a cabeça da cobra
	var newHead point

	// Calcula a próxima posição da cabeça da cobra com base na direção
	switch direction {
	case up:
		newHead = point{head.x, head.y - 1}
	case down:
		newHead = point{head.x, head.y + 1}
	case left:
		newHead = point{head.x - 1, head.y}
	case right:
		newHead = point{head.x + 1, head.y}
	}

	// Verifica se a nova cabeça colide com a cobra ou com as paredes
	for _, segment := range snake {
		if newHead == segment || newHead.x < 1 || newHead.x >= width-1 || newHead.y < 1 || newHead.y >= height-1 {
			gameOver = true // Define a flag de fim de jogo como verdadeira
			return
		}
	}

	// Verifica se a nova cabeça está na posição da comida
	if newHead == food {
		score++     // Aumenta a pontuação
		placeFood() // Coloca a comida em uma nova posição aleatória
	} else {
		// Se não estiver na comida, remove o último segmento da cobra
		snake = snake[:len(snake)-1]
	}

	// Adiciona a nova cabeça à cobra
	snake = append([]point{newHead}, snake...)
}

// Função placeFood coloca a comida em uma posição aleatória na tela
func placeFood() {
	rand.Seed(time.Now().UnixNano())                              // Inicializa a semente do gerador de números aleatórios
	food = point{rand.Intn(width-2) + 1, rand.Intn(height-2) + 1} // Posição aleatória dentro das bordas
}

// Função render desenha a tela do jogo
func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault) // Limpa a tela

	// Desenha as bordas
	for x := 0; x < width; x++ {
		termbox.SetCell(x, 0, borderChar, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, height-1, borderChar, termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < height; y++ {
		termbox.SetCell(0, y, borderChar, termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(width-1, y, borderChar, termbox.ColorWhite, termbox.ColorDefault)
	}

	// Desenha a cobra
	for _, segment := range snake {
		termbox.SetCell(segment.x, segment.y, '█', termbox.ColorGreen, termbox.ColorDefault)
	}

	// Desenha a comida
	termbox.SetCell(food.x, food.y, '■', termbox.ColorRed, termbox.ColorDefault)

	// Desenha a pontuação
	scoreStr := fmt.Sprintf("Score: %d", score)
	for i, char := range scoreStr {
		termbox.SetCell(i, height, char, termbox.ColorWhite, termbox.ColorDefault)
	}

	if gameOver {
		gameOverMessage1 := []rune("Game Over! Press Esc to Exit.")
		gameOverMessage2 := []rune("Or")
		gameOverMessage3 := []rune("Enter to Restart Game")
		gameOverMessage4 := []rune(fmt.Sprintf("Your Score: %d", score))
		for i, char := range gameOverMessage1 {
			termbox.SetCell((width-len(gameOverMessage1))/2+i, height/2-1, char, termbox.ColorWhite, termbox.ColorDefault)
		}
		for i, char := range gameOverMessage2 {
			termbox.SetCell((width-len(gameOverMessage2))/2+i, height/2, char, termbox.ColorWhite, termbox.ColorDefault)
		}
		for i, char := range gameOverMessage3 {
			termbox.SetCell((width-len(gameOverMessage3))/2+i, height/2+1, char, termbox.ColorWhite, termbox.ColorDefault)
		}
		for i, char := range gameOverMessage4 {
			termbox.SetCell((width-len(gameOverMessage4))/2+i, height/2+2, char, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.Flush() // Atualiza a tela
}
