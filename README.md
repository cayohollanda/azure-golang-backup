# Azure Golang Backup Script
Um script feito em Golang com o objetivo de automatizar o processo de backup de dados na Azure Storage

## Uso
Primeiro devemos adicionar o executável ao PATH do sistema, usando: 
```console
cayohollanda@pc:~$ sudo cp azure-backup /usr/local/bin/
```
Após isso, precisamos adicionar os seguintes environments: 
* <b>AZURE_STORAGE_ACCOUNT</b>: Environment onde fica definido o nome da conta
* <b>AZURE_STORAGE_ACCESS_KEY</b>: Environment que fica alocado a access key de autenticação
* <b>AZURE_STORAGE_CONTAINER_NAME</b>: Environment que guarda o nome do container para onde vai o blob
* <b>AZURE_DIRECTORY_TO_UPLOAD</b>: Environment onde fica o diretório que será salvo na Azure

Após setarmos os environments, rodamos o nosso script: 
```console
cayohollanda@pc:~$ azure-backup
```

E pronto, basta aguardar o script fazer a compressão dos arquivos e o upload para a Azure.

## Agendamentos e automatização
Para agendamentos e automatizaço do uso do script, você pode seguir o rápido tutorial que é passado na <a href="https://github.com/cayohollanda/aws-golang-backup">versão AWS deste script</a>, ao fim do README, 
é passada uma rápida explicação do uso do script com o Cron.
