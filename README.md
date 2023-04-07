# cardinal_chains

#Project Description
Cardinal Chains is an open-source puzzle game that challenges players to create chains of increasing integers by connecting adjacent cells in a grid. The objective of the game is to form the longest possible chain, starting from the lowest number and ending with the highest number, while following specific game rules.

#Installing and Running Cardinal Chains
In this guide, we will walk you through the steps to install and run the Cardinal Chains project on your machine. We will be using Windows Subsystem for Linux 2 (WSL 2) and installing the libyaml library.

##Step 1: Install WSL 2
Follow the official Microsoft guide to set up WSL 2 on your Windows machine:
https://docs.microsoft.com/en-us/windows/wsl/install

Make sure to complete all the steps, including updating to WSL 2, installing a Linux distribution (such as Ubuntu), and setting up your new Linux distribution.

##Step 2: Install libyaml
Once you have WSL 2 set up and your Linux distribution installed, open a new WSL terminal.

Update your package list and install the required dependencies:
```bash
sudo apt update
sudo apt install -y libyaml-dev
```

##Step 3: Clone the Cardinal Chains repository
In the WSL terminal, navigate to the directory where you want to store the project and clone the GitHub repository:
```bash
git clone https://github.com/Xefreh/cardinal_chains.git
cd cardinal_chains
```

##Step 4: Install CMake
Before you can build the project, you need to have CMake installed. Run the following commands to install CMake:
```bash
sudo apt update
sudo apt install -y cmake
```

##Step 5: Build the project using CMake
Create a build directory inside the project folder and navigate to it:
```bash
mkdir build
cd build
```

Now run CMake to generate the Makefiles and build the project:
```bash
cmake ..
make
```

##Step 6: Run Cardinal Chains
After the build is complete, you can run Cardinal Chains using the following command:
```bash
./cardinal_chains ../levels.yml
```
This command will execute the cardinal_chains binary, loading the game levels from the levels.yml file.

Now you can enjoy playing Cardinal Chains!
