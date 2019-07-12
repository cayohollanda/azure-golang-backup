package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/cayohollanda/azure-golang-backup/utils"
)

func main() {
	utils.TimedPrintln("Iniciando rotina...")

	var (
		azureAccount   = os.Getenv("AZURE_STORAGE_ACCOUNT")
		azureAccessKey = os.Getenv("AZURE_STORAGE_ACCESS_KEY")
		containerName  = os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
		uploadPtr      = flag.Bool("u", true, "Valor booleano que define se irá fazer upload do arquivo")
		// zipPtr                = flag.Bool("z", true, "Valor booleano que define se irá ou não zipar o arquivo")
		downloadPtr           = flag.String("d", "", "Valor em texto que define qual arquivo será feito o download")
		directoryToUploadPtr  = flag.String("directory", "", "Valor em texto que define qual o diretório de upload")
		azureAccountNamePtr   = flag.String("azureAccountName", "", "Account name da Azure")
		azureAccessKeyPtr     = flag.String("azureAccessKey", "", "Access Key de acesso à conta de armazenamento da Azure")
		azureContainerNamePtr = flag.String("azureContainerName", "", "Container name da conta de armazenamento da Azure")
		uploadValue           bool
		downloadValue         string
		// zipValue                    bool
		directoryToUploadValue      string
		azureAccountParameter       string
		azureAccessKeyParameter     string
		azureContainerNameParameter string
	)

	flag.Parse()

	uploadValue = *uploadPtr
	downloadValue = *downloadPtr
	// zipValue = *zipPtr
	directoryToUploadValue = *directoryToUploadPtr
	azureAccountParameter = *azureAccountNamePtr
	azureAccessKeyParameter = *azureAccessKeyPtr
	azureContainerNameParameter = *azureContainerNamePtr

	if len(azureAccount) == 0 || len(azureAccessKey) == 0 || len(containerName) == 0 {
		if len(azureAccountParameter) == 0 || len(azureAccessKeyParameter) == 0 || len(azureContainerNameParameter) == 0 {
			utils.TimedPrintln("Defina as variáveis de ambiente: AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_ACCESS_KEY e AZURE_STORAGE_CONTAINER_NAME")
			os.Exit(1)
		} else {
			azureAccount = azureAccountParameter
			azureAccessKey = azureAccessKeyParameter
			containerName = azureContainerNameParameter
		}
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

	if uploadValue {
		// filename := ""
		// blobURL := containerURL.NewBlockBlobURL(filename)

		blobsList, err := containerURL.ListBlobsFlatSegment(context, azblob.Marker{}, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Println(err.Error())
		}

		if len(blobsList.Segment.BlobItems) > 0 {
			newBlob := blobsList.Segment.BlobItems[0]
			for _, blob := range blobsList.Segment.BlobItems {
				if blob.Properties.LastModified.String() > newBlob.Properties.LastModified.String() {
					newBlob = blob
				}
			}

			log.Println(newBlob.Name)
		}

		// if zipValue {
		// 	utils.TimedPrintln("Zipando arquivo...")

		// 	var tagDate string
		// 	t := time.Now()
		// 	month, _ := strconv.Atoi(fmt.Sprintf("%d", t.Month()))
		// 	day, _ := strconv.Atoi(fmt.Sprintf("%d", t.Day()))
		// 	if day < 10 {
		// 		if month < 10 {
		// 			tagDate = fmt.Sprintf("0%d0%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
		// 		} else {
		// 			tagDate = fmt.Sprintf("0%d%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
		// 		}
		// 	} else {
		// 		if month < 10 {
		// 			tagDate = fmt.Sprintf("%d0%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
		// 		} else {
		// 			tagDate = fmt.Sprintf("%d%d%d%02d%02d%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
		// 		}
		// 	}

		// 	filename = fmt.Sprintf("blob%s", tagDate)
		// 	utils.ZipWriter(directoryToUploadValue, filename) // criando arquivo zip
		// 	blobURL = containerURL.NewBlockBlobURL(filename)
		// } else {
		// 	filename = directoryToUploadValue
		// }

		// file, err := os.Open(filename)
		// utils.CheckErr("Erro ao verificar se o arquivo existe", err)

		// log.Println("Fazendo upload na Azure")
		// _, err = azblob.UploadFileToBlockBlob(context, file, blobURL, azblob.UploadToBlockBlobOptions{})
		// utils.CheckErr("Erro ao fazer upload do blob", err)
		// utils.TimedPrintln("Upload feito com sucesso!")
		// utils.TimedPrintln("Backup feito com o nome: " + filename)

		// _ = os.Remove(filename)
	}

	if downloadValue != "" {
		utils.TimedPrintln("Fazendo download do blob...")
		blobURL := containerURL.NewBlockBlobURL(downloadValue)
		_, err = blobURL.Download(context, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
		utils.CheckErr("", err)
		utils.TimedPrintln("Download feito com sucesso")
	}
}
