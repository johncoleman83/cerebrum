#!/usr/bin/env bash

# npm install -g multi-file-swagger

function introduction_message() {

	go get -t -v ./...
}

function special_execute() {
    set -x; "$@"; set +x;
}

function check_if_should_execute_command() {
    while true; do
        read -p "Do you wish to execute '$1'? " -n 1 -r yn
        case $yn in
            [Yy]* ) special_execute "$1"; break;;
            [Nn]* ) echo -e "\ncontinue without executing '$1'\n"; break;;
            * ) -e echo "\nPlease answer Yes or No.\n";;
        esac
    done
}

function continue_loop() {
    while true; do
        read -p "Do you wish to continue? " -n 1 -r yn
        case $yn in
            [Yy]* ) break;;
            # exit if for some reason, the path is empty string
            [Nn]* ) exit 1;;
            * ) echo -e "\nPlease answer Yes or No.\n";;
        esac
    done
}

function execute_setup() {
    continue_loop

    echo "This app uses make"
    check_if_should_execute_command "make help"

    echo "Node is a dependency for the swagger docs"
    check_if_should_execute_command "brew install node"
    check_if_should_execute_command "npm install -g multi-file-swagger"

    echo "we use some go support executables"
    check_if_should_execute_command "go get -u golang.org/x/lint/golint"
    check_if_should_execute_command "brew install dep"

    echo "do you want go to fetch & install all the dependencies?"
    check_if_should_execute_command "go get -t -v ./..."
}

execute_setup
