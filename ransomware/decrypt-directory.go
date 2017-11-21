package main
import "crypto/aes"
import "crypto/cipher"
import "io/ioutil"
import "io"
import "os"
import "bytes"
import "log"
import "strings"
import "fmt"
import "bufio"

var encryption_extension string = "__enc__"

// This function decrypts the directory and its children in a recursive way
func decryptDirectory(nameDir string, key []byte) {
    files, err := ioutil.ReadDir(nameDir)
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
        // This will be the absolute path to either the directory or the path
        var fullChildPath string = nameDir + "/" + f.Name()

        // We check the file type of this file/directory
        fileInfo, err := os.Stat(fullChildPath)
        if err != nil {
            panic(err.Error())
        }
        // If it's a directory, we transverse recursively and concurrently
        if fileInfo.IsDir() {
            decryptDirectory(fullChildPath, key)
        } else {
            // We decrypt if it's a file that contains the encryption extension
            if strings.Contains(f.Name(), encryption_extension){
                decryptFile(fullChildPath, key)
            }
            
        }
    }
}
func decryptFile(name string, key []byte) {
    ciphertext, err := ioutil.ReadFile(name)
    if err != nil {
        panic(err.Error())
    }
    
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    // Get the 16 byte IV
    iv := ciphertext[:aes.BlockSize]

    // Remove the IV from the ciphertext
    ciphertext = ciphertext[aes.BlockSize:]

    // Create a stream for decrypting the ciphertext
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)
    // Remove the encrypted file
    os.Remove(name)
    
    // We get the name of the file prior to its encryption
    name = strings.TrimSuffix(name,encryption_extension)
    // Create file with the decrypted information
    createFileWithContents(name, ciphertext) 
}

func createFileWithContents(file  string, text []byte){
    // We create a new file for saving the encrypted data.
    f, _ := os.Create(file)
    _, _ = io.Copy(f, bytes.NewReader(text))
 
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter decryption key: ")
    key_16, _ := reader.ReadString('\n')
    // We remove the new line at the end of key_16
    key_16 = key_16[:len(key_16) - 1]
    fmt.Println(key_16)
    dir, _ := os.Getwd()
    key := []byte(key_16)
    decryptDirectory(dir,key)    
}
