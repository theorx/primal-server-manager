## PSM for running rust servers
*Author: Lauri Orgla*




----
# TODO:

### Design a scheduling rule engine
* Design event sourcing system which would be able to determine 1) oxide update, 2) whether it was a facepunch monthly update wipe or not

* Implement a recurring rule engine based on
* What happens on

### Server definitions
* Everything related to specifying a server

### Plugin definitions
* Plugin has a source file that is mutable, ex: TickrateLimiter 
    * Plugin name
    * Plugin configuration files (Versioned based on server) - multi config support / generated config or something based on server definitions vars
    * Have the ability to list all of the "keys" used in config for dynamic generation based on "server definition vars"

      //Have a mechanism to define "plugins"

      //Plugins are c# sources that can have files associated with them
      // 		- The file associations are of two kind:
  