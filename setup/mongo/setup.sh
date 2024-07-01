#!/bin/bash

sudo mkdir -p $HOME/mongodb/mchosting/db

# Setting the owner to the currently logged-in user
sudo chown -R $(whoami) $HOME/mongodb/mchosting/db

# Set appropriate permissions
sudo chmod 750 $HOME/mongodb/mchosting/db
sudo chmod 750 $HOME/mongodb/mchosting

sudo mongod --config config.conf --fork

sleep 10

mongosh setup.js

sudo mongod --config config.conf --auth --fork