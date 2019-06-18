package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/briandowns/spinner"

	"github.com/cayohollanda/azure-golang-backup/utils"
)

func main() {
	utils.TimedPrintln("Iniciando rotina...")

	var (
		azureAccount           = os.Getenv("AZURE_STORAGE_ACCOUNT")
		azureAccessKey         = os.Getenv("AZURE_STORAGE_ACCESS_KEY")
		containerName          = os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
		uploadPtr              = flag.Bool("u", true, "Valor booleano que define se irá fazer upload do arquivo")
		zipPtr                 = flag.Bool("z", true, "Valor booleano que define se irá ou não zipar o arquivo")
		downloadPtr            = flag.String("d", "", "Valor em texto que define qual arquivo será feito o download")
		directoryToUploadPtr   = flag.String("directory", "", "Valor em texto que define qual o diretório de upload")
		uploadValue            bool
		downloadValue          string
		zipValue               bool
		directoryToUploadValue string
	)

	flag.Parse()

	uploadValue = *uploadPtr
	downloadValue = *downloadPtr
	zipValue = *zipPtr
	directoryToUploadValue = *directoryToUploadPtr

	if len(azureAccount) == 0 || len(azureAccessKey) == 0 || len(containerName) == 0 {
		utils.TimedPrintln("Defina as variáveis de ambiente: AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY e AZURE_STORAGE_CONTAINER_NAME")
		os.Exit(1)
	}

	if directoryToUploadValue == "" && uploadValue && downloadValue == "" {
		utils.TimedPrintln("Defina o diretório para upload usando -directory=<diretorio>")
		os.Exit(1)
	}

	utils.TimedPrintln("Setando as credenciais...")
	credentials, err := azblob.NewSharedKeyCredential(azureAccount, azureAccessKey)
	utils.CheckErr("Erro com as credenciais informadas", err)
	utils.TimedPrintln("Credenciais setadas")

	pipeline := azblob.NewPipeline(credentials, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: 5 * time.Minute,
		},
	})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", azureAccount, containerName))

	containerURL := azblob.NewContainerURL(*URL, pipeline)

	context := context.Background()

	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)

	if uploadValue {
		filename := ""
		blobURL := containerURL.NewBlockBlobURL(filename)

		if zipValue {
			utils.TimedPrintln("Zipando arquivo...")

			var tagDate string
			t := time.Now()
			month, _ := strconv.Atoi(fmt.Sprintf("%d", t.Month()))
			day, _ := strconv.Atoi(fmt.Sprintf("%d", t.Day()))
			if day < 10 {
				if month < 10 {
					tagDate = fmt.Sprintf("0%d0%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
				} else {
					tagDate = fmt.Sprintf("0%d%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
				}
			} else {
				if month < 10 {
					tagDate = fmt.Sprintf("%d0%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
				} else {
					tagDate = fmt.Sprintf("%d%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
				}
			}

			filename = fmt.Sprintf("blob%s", tagDate)
			utils.ZipWriter(directoryToUploadValue, filename) // criando arquivo zip
			blobURL = containerURL.NewBlockBlobURL(filename)
		} else {
			filename = directoryToUploadValue
		}

		file, err := os.Open(filename)
		utils.CheckErr("Erro ao verificar se o arquivo existe", err)

		s.Prefix = "Fazendo upload na Azure"
		s.Color("red", "bold")
		s.Start()
		_, err = azblob.UploadFileToBlockBlob(context, file, blobURL, azblob.UploadToBlockBlobOptions{})
		utils.CheckErr("Erro ao fazer upload do blob", err)
		utils.TimedPrintln("Upload feito com sucesso!")
		s.Stop()
		utils.TimedPrintln("Backup feito com o nome: " + filename)

		_ = os.Remove(filename)
	}

	if downloadValue != "" {
		utils.TimedPrintln("Fazendo download do blob...")
		blobURL := containerURL.NewBlockBlobURL(downloadValue)
		_, err = blobURL.Download(context, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
		utils.CheckErr("", err)
		utils.TimedPrintln("Download feito com sucesso")
	}
}
