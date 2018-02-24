#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

###############################################################################
# MAIN
printf "The kubed-sh end-to-end test can take up to 2 min.\n"

printf "TC: one-shot dproc (binary)\n"
./dproc-one-shot.kbdsh

printf "TC: long-running dproc (Python)\n"
./dproc-longrunning.kbdsh

"\nDONE=====================================================================\n"
