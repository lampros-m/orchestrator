# Delete all .log files in the current directory and subdirectories forcefully
find . -type f -name "*.log" -exec rm -f {} \;