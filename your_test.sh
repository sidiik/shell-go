#!/bin/bash
# test.sh - Example shell script

# A simple greeting function
greet() {
    local name=$1
    echo "Hello, $name!"
}

# Check if a name argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <name>"
    exit 1
fi

# Call the greeting function
greet "$1"

# Loop example
echo "Counting from 1 to 5:"
for i in {1..5}; do
    echo "Number $i"
done

# Conditional example
echo "Checking if a file named 'example.txt' exists..."
if [ -f "example.txt" ]; then
    echo "File exists!"
else
    echo "File does not exist."
fi

# End of script
echo "Script finished."
