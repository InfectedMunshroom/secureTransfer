# Things to do:

1. There is a file known as upload.go in folder client->upload
2. The main.go file in the root directory runs the server that will allow upload of the files. The configuration of the server is present in folder internal->server->server.go
3. Make a GUI that calls the UploadFileWithAES taking the following inputs from the user: filename and url of destination server
   1. UploadFileWithAES takes in 4 inputs:
      1. aeskey := provide it yourself for the time being
      2. filename := user provides it
      3. destination url := you can hardcode it or the user can provide it, it's for the server that will take the file. 
         1. syntax: http://localhost:8081/upload
      4. rsaFile := just put in the path for final.pub
4. Ensure to run the main.go file on a different port, **do not merge it with the server that will run the GUI page**
   1. For eg. if you are running the GUI server on localhost:8080, run the upload file server (main.go) on localhost:8081 or something


# Demo Dry Run:
1. Server for GUI starts on localhost:8080
2. Server for accepting files starts on localhost:8081 by running main.go
3. User inputs filepath required for upload in gui
4. GUI takes the filepath, and then calls the UploadFileWithAES function from upload.go
5. In the directory where main.go exists (aka the root directory) 2 folders aka aes and file_images should be created if everything runs successfully

# Main dry run:
1. User encrypts text into file using steganography
2. User then selects "upload" which takes the absolute file path of the now encrypted .png file and calls the UploadFileWithAES to upload it to the server
3. Everything goes as planned and everyone lives happily ever after


# Testing to do before starting gui
1. Run main.go on localhost:8080
2. Make a client.go file in the secureTransfer folder, put package as main
3. Call UploadFileWithAES from encryptdecrypt package, execute it and ensure the files get uploaded with the corrected messages coming on both terminals
4. Voila, now you can work on creating the gui with full understanding of everything you need to understand
 