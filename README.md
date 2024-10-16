# Additonal Features
1. Doctor, Patient database
2. Different platforms :-: portable
3. RBAC 
4. Key Rotation :-: Cipher Rotation
5. Attestations -> Blockchain -> present on all nodes


# Things to do

## Stage 1: Finish Managed File Transfer System Software
1. RBAC with Server-based client login
2. Attestations
   1. Try to put them on blockchain if possible (present on all nodes)
3. Different platforms
   1. Render the front end for downloading files available on Windows and Linux (primary), Android (secondary)
   2. See if we can make the encryption process for android
4. Key & Cipher Rotation
   1. Implement real-time random creation of AES keys (primary)
   2. Allow retrieval of public keys through channels, secure private keys
   3. Allow possibility of cipher rotation (Not really needed but sure)
5. Refactor, deletion of temporary files and organisation of file structure in server

Priority Order: 1 -> 4.1 -> 2 -> 5 -> 4.2 -> 3 -> 4.3

## Stage 2: Customisation of the project to emulate it's usage in a hospital
1. I guess we should fork it and then do it, leave this untouched after stage 1 is done
