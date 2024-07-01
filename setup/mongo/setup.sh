#!/bin/bash

mkdir -p $HOME/mongodb/mchosting/db

# Setting the owner to the currently logged-in user
CUR_USER="$(whoami)"

sudo chown -R $CUR_USER $HOME/mongodb/mchosting/db

# Set appropriate permissions
sudo chmod 750 $HOME/mongodb/mchosting/db
sudo chmod 750 $HOME/mongodb/mchosting

sudo mongod --config config.conf --fork

sleep 10

mongosh --port 27777 setup.js

sudo mongod --config config.conf --auth --fork