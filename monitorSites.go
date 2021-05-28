package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 5
const delay = 10
const arquivoLog = "log.txt"
const arquivoSites = "sites.txt"

func main() {
	exibeIntroducao()
	for {
		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			exibeLogs()
		case 3:
			limparLogs()
		case 0:
			os.Exit(0)
		default:
			fmt.Println("Nao reconheco este comando!")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	nome := "Rafael"
	versao := 1.1
	fmt.Println("Ola, sr.", nome)
	fmt.Println("Este programa esta na versao ", versao)
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	return comandoLido
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir os Logs")
	fmt.Println("3 - Limpar Logs")
	fmt.Println("0 - Sair do programa")
	fmt.Println("")
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := lerSitesDoArquivo()
	for i := 0; i < monitoramentos; i++ {
		for _, site := range sites {
			fmt.Println("Testando site: ", site)
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
		registraLog(site, false)
	}
	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, " foi carregado com sucesso!")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, " esta com problemas. Status Code: ", resp.StatusCode)
		registraLog(site, false)
	}
}

func lerSitesDoArquivo() []string {

	var sites []string

	arquivo, err := os.Open(arquivoSites)
	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		sites = append(sites, strings.TrimSpace(linha))
		if err == io.EOF {
			break
		}
	}
	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile(arquivoLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - Online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func exibeLogs() {
	arquivo, err := ioutil.ReadFile(arquivoLog)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	fmt.Println(string(arquivo))
}

func limparLogs() {
	os.Remove(arquivoLog)
}
