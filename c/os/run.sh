set -e
./build.sh
bochs -f bochsrc.txt -q -rc debug.rc 
