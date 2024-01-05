# Cardinal Chains

Cardinal Chains is a puzzle game that challenges players to create sequences by connecting numbers while following specific rules.

## Prerequisites

Before you begin, ensure you have met the following requirements:
- You have a Linux or macOS machine. (Windows steps may vary)
- You have installed CMake.
- You have a C compiler installed, such as GCC.
- You have installed the LibYAML library.

## Installing Cardinal Chains

To install Cardinal Chains, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/Xefreh/cardinal_chains.git
   ```
2. Navigate to the project directory:
   ```
   cd cardinal_chains
   ```

## Compiling the Project

Cardinal Chains uses CMake for its build system. To compile the project, follow these steps:

1. Create and navigate to the build directory:
   ```
   mkdir build && cd build
   ```
2. Generate the build system files with CMake:
   ```
   cmake ..
   ```
3. Compile the project:
   ```
   make
   ```

## Running the Program

After compiling the project, you can run the Cardinal Chains game from the build directory:

```
./cardinal_chains
```

