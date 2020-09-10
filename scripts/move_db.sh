#!/bin/sh

DATA_DIRECTORY=$HOME/.local/share/qotd

if [ ! -d "$DATA_DIRECTORY" ]; then
    mkdir $DATA_DIRECTORY
fi
cp quran.db $DATA_DIRECTORY/quran.db