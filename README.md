goPlay
======

Built to interact with hitbox.tv chat during streams specifically during games. Code is hacked together but works.

Status
-----

Basically done for what I plan to do with it. I found the xdotool to be finicky with mednafen but eventually got it to work.

Building/Running
-----
Build instructions are for Ubuntu Server 12.04 LTS. They should work out of the box for other Ubuntu versions and can be easily tweaked for other distros. 

Go is assumed to have been installed.

Windows is completely unsupported but running in virtualbox is trivial.

    sudo aptitude install mednafen
    sudo aptitude install xdotool
    go get github.com/gorilla/websocket
    //cfg allows you to work with sane keybindings
    cp mednafen.cfg !mednanfen install directory!
    open your game of choice in mednafen
    nano interactive.go //change channel="" located in main() to channel to monitor
    go run interactive.go

You can build to send random input to the generator via the commented section at the bottom. Be warned that progress is nonexistent without serious rng biasing.

A gnuboy config file is also included if you wish to interface with that specific emulator via xdotool.
