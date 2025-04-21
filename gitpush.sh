#!/bin/bash

# Commit all changes with a default message or your custom message
git add .
git commit -m "${1:-Auto update}"
git push origin master


