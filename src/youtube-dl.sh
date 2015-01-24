#!/bin/bash
cd $3;
/usr/local/bin/youtube-dl -x --audio-format mp3 $1 --exec "cp {} $2 && rm {}";