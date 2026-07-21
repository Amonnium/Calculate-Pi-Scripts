#!/bin/bash
export BC_LINE_LENGTH=0
bc -lq <<< "scale=1000; 4*a(1)"   