# Network-Project

## Stateless password manager

 - A password manager is designed to help a user store and manage online credentials. 
 - Stateless password managers do not require to store a database since the passwords are generated randomly using the master password and a key generation function.
 - The project aims to create a stateless password manager service. 
    - The client can enter the name of the website and a strong password is automatically generated. 
    - To retrieve the password, the process is the same

### Features

 - Super Secure Encryption 
 - Stateless password server: Does not need database.
 - One master password for all your passwords :P
 - Command line interface - Ease of use
 - High speed service - Developed in GoLang


### Working 

#### Client Side 

 - /getPass
    - generates password given username, master secret and website name
 - /help
    - displays help message
 - /getSpec
     - prints out the specifications of the program
 - /quit
    - quit server connection

#### Server Side

 - Listen to client through the socket. 
    - Module: net.Listen(, “ip:port”)
 - Accept packet from Client using a TCP socket connection.
 - Use master password to generate new password for the server.
 - Modules
     - Generate Password (‘>/gen_pass’) 
     - Quit (>/quit)
 - Generates secure password from the genPASS(username,website, secret)

