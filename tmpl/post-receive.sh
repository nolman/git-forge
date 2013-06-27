#!/bin/bash

RECEIVE_REPO={{.AppName}} RECEIVE_USER=git DATABASE_URL={{.DatabaseUrl}} gitreceive hook
