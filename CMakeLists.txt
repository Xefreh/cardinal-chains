cmake_minimum_required(VERSION 3.0)
project(cardinal_chains)

# Set the LibYAML library directory manually
set(LIBYAML_LIBRARY_DIRS /usr/lib/x86_64-linux-gnu/)

# Find the LibYAML library
find_library(LIBYAML_LIBRARIES NAMES yaml HINTS ${LIBYAML_LIBRARY_DIRS})

# Include the LibYAML directories
include_directories(/usr/lib/x86_64-linux-gnu/)

# Create an executable target
add_executable(cardinal_chains main.c)

# Link against the LibYAML library
target_link_libraries(cardinal_chains ${LIBYAML_LIBRARIES})

# Copy the input file to the build directory
configure_file(levels.yml ${CMAKE_CURRENT_BINARY_DIR}/levels.yml COPYONLY)

# Set the input file path as a preprocessor definition
target_compile_definitions(cardinal_chains PRIVATE INPUT_FILE="${CMAKE_CURRENT_BINARY_DIR}/levels.yml")