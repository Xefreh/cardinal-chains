#include "input.h"

#include <stdio.h>

char get_user_input(void) {
    printf("Select a direction (N, S, E, W).\n");
    printf("Cancel the previous move (B).\n");
    printf("Erase the chain (R).\n");
    printf("Restart the level (X).\n");
    printf("Select another chain (C).\n");
    return getchar();
}
