#!/bin/bash
# Get the effective user ID
euid=$(id -u)

# Check if eUID is 0 (root)
if [ $euid -eq 0 ]; then
  echo "Script executed with sudo permissions."
else
  echo "Script does not have sudo permissions. Please rerun with sudo."
  exit 1  
fi

# Define an array of dependencies (replace with actual package names)
dependencies=("go" "make" "ip" "aircrack-ng")

# Function to check if a package is installed
function is_installed() {
  package_name="$1"
  if command -v "$package_name" >/dev/null 2>&1; then
    return 0  # package is installed
  else
    return 1  # package is not installed
  fi
}

# Check for dependencies
for dependency in "${dependencies[@]}"; do
  if ! is_installed "$dependency"; then
    echo "Error: Dependency '$dependency' is not installed."
    echo "Please install it and try again."
    exit 1
  fi
done

# Get the program name from the Makefile (optional)
program_name=raspberry  # Or use $1 as an argument

# Run the make script in the folder
make "$program_name"

if [ $? -eq 0 ]; then
  echo "Successfully compiled $program_name!"
else
  echo "Error compiling $program_name. Please check the logs for details."
  exit 1
fi

# Install directory (replace with your desired location)
install_dir="/usr/local/bin"

# Check if the program is built
if [ ! -f "./bin/$program_name" ]; then
  echo "Error: Program '$program_name' not found. Please build it first."
  exit 1
fi

# Copy the program to the installation directory
sudo cp "./bin/$program_name" "$install_dir"
if [ $? -eq 0 ]; then
  echo "Program '$program_name' copied successfully."
else
  echo "Error copying program. Please check permissions."
  exit 1
fi

# Make the directory in the homefolder
mkdir -p "/usr/local/raspberry"

sudo chgrp sudo /usr/local/raspberry
sudo chmod g+rwx /usr/local/raspberry

# Check if the directory creation was successful
if [ $? -eq 0 ]; then
  echo "Directory '/usr/local/raspberry' created successfully."
else
  echo "Error creating directory. Please check permissions."
  exit 1  # Exit script with error code
fi

# Copy the source directory to the new directory
cp -r "./ui" "/usr/local/raspberry"

# Check if the copy operation was successful
if [ $? -eq 0 ]; then
  echo "Interface copied successfully."
else
  echo "Error copying directory. Please check permissions or source directory existence."
  exit 1  # Exit script with error code
fi

# Add the install directory to PATH (optional - modify with caution)
# This line is commented out for safety reasons
# echo "export PATH=$PATH:$install_dir" >> ~/.bashrc

echo "**Note:** You may need to add '$install_dir' to your PATH environment variable manually for the program to be accessible from any location."

