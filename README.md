# TCP Chat app

This is a very simple chat app using TCP.  A user will broadcast their messages to all other people connected.

## Added the following commands for the client
* JOIN - You must "join" before you can start posting messages.
    * join (username) [password required for 'admin']
    * join (username) -replace {__still needs to be implemented__}
* HELP - Display the list of commands
* WHOAMI - Display your username..
* DEBUG - This is only allowed for Admin user.


## TODO 
* Don't allow the same username to be entered on two sessions.
* Store Admin password in an environment variable rather than being hard-coded.