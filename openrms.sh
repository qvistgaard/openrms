#!/usr/bin/env bash
git pull
go build github.com/qvistgaard/openrms/cmd/openrms
gnome-terminal --full-screen --profile OpenRMS --hide-menubar --wait -- ./openrms -driver oxigen