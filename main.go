package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"./utils"
)

func main() {
	utils.TimedPrintln("Iniciando rotina...")

	var (
		azureAccount      = os.Getenv("AZURE_STORAGE_ACCOUNT")
		azureAccessKey    = os.Getenv("AZURE_STORAGE_ACCESS_KEY")
		containerName     = os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
		directoryToUpload = os.Getenv("AZURE_DIRECTORY_TO_UPLOAD")
	)

	if len(azureAccount) == 0 || len(azureAccessKey) == 0 || len(containerName) == 0 || len(directoryToUpload) == 0 {
		utils.TimedPrintln("Defina as variáveis de ambiente: AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY e AZURE_STORAGE_CONTAINER_NAME, AZURE_DIRECTORY_TO_UPLOAD")
		os.Exit(1)
	}

	utils.TimedPrintln("Setando as credenciaiss...")
	credentials, err := azblob.NewSharedKeyCredential(azureAccount, azureAccessKey)
	utils.CheckErr("Erro com as credenciais informadas", err)
	utils.TimedPrintln("Credenciais setadas")

	pipeline := azblob.NewPipeline(credentials, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", azureAccount, containerName))

	containerURL := azblob.NewContainerURL(*URL, pipeline)

	context := context.Background()

	utils.TimedPrintln("Zipando arquivo...")
	timeBlob := time.Now()
	zipFilename := fmt.Sprintf("blob%d%d%d%02d%02d%02d", timeBlob.Day(), timeBlob.Month(), timeBlob.Year(), timeBlob.Hour(), timeBlob.Minute(), timeBlob.Second())
	utils.ZipWriter(directoryToUpload, zipFilename) // criando arquivo zip

	blobURL := containerURL.NewBlockBlobURL(zipFilename)
	file, err := os.Open(zipFilename)
	utils.CheckErr("Erro ao verificar se o arquivo zip existe", err)

	utils.TimedPrintln("Fazendo upload do blob para a Azure...")
	_, err = azblob.UploadFileToBlockBlob(context, file, blobURL, azblob.UploadToBlockBlobOptions{})
	utils.CheckErr("Erro ao fazer upload do blob", err)
	utils.TimedPrintln("Upload feito com sucesso!")

	utils.TimedPrintln("Fazendo download do blob...")
	_, err = blobURL.Download(context, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	utils.CheckErr("", err)
	utils.TimedPrintln("Download feito com sucesso")

	os.Remove(zipFilename)

	utils.TimedPrintln("Backup feito com o nome: " + zipFilename)
}
