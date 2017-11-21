# Ransomware
This crypto ransomware is for academic and learning purposes. It encrypts every folder in the directory where it is ran. It also encrypts files found in subdirectories.

## How to run

This ransomware is divided in two parts: the PHP server that receives API calls for creating/retrieving encryption keys, and the cryptoware-ransomware itself, which is a Go app. 

### Server

1. Place the files inside server/ransomware inside a folder named "ransomware". 
2. Move this folder to your PHP server. 

### Ransomware

You need Go for creating the executable binary of the ransomware. 

The ransomware consists of the encrypter and the decrypter. 

#### Building for Intel based computers running Windows:
  
  Building the decrypter
  ```
    env GOOS=windows GOARCH=386 go build decrypt-directory.go
  ```
  Building the encrypter
  ```
    env GOOS=windows GOARCH=386 go build encrypt-directory.go
  ```

#### Building for other OS's and architectures? 
Take a look at [this](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04) guide:

## Architecture

Segmented arrows stand for asynchronous calls.

![architecture](https://github.com/edfullest/Ransomware/blob/master/architecture.png)
