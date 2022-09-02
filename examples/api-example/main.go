package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ErrMsg struct{ error }

var help = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type Model struct {
	message string
}

func (m Model) Init() tea.Cmd {
	return getJoke
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Joke:
		m.message = msg.Joke
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, getJoke
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.message, help("q: quit • n: next joke"))
}

type Joke struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Joke     string `json:"joke"`
	Flags    struct {
		Nsfw      bool `json:"nsfw"`
		Religious bool `json:"religious"`
		Political bool `json:"political"`
		Racist    bool `json:"racist"`
		Sexist    bool `json:"sexist"`
		Explicit  bool `json:"explicit"`
	} `json:"flags"`
	Safe bool   `json:"safe"`
	ID   int    `json:"id"`
	Lang string `json:"lang"`
}

func getJoke() tea.Msg {
	resp, err := http.Get("https://v2.jokeapi.dev/joke/Programming?blacklistFlags=nsfw,religious,political,racist,sexist,explicit&type=single")
	if err != nil {
		return ErrMsg{fmt.Errorf("unable to create get request: %w", err)}
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrMsg{fmt.Errorf("unable to read response: %w", err)}
	}
	var joke Joke
	err = json.Unmarshal(data, &joke)
	if err != nil {
		return ErrMsg{fmt.Errorf("unable to unmarshal json: %w", err)}
	}
	return joke
}

func main() {
	m := Model{}
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
