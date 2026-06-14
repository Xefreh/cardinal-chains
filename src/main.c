#include <stdio.h>
#include <stdlib.h>

#include "game.h"
#include "game_loop.h"
#include "level.h"
#include "level_loader.h"

int main(int argc, char *argv[]) {
    if (argc < 2) {
        fprintf(stderr, "Usage: %s <YAML file>\n", argv[0]);
        return 1;
    }

    Levels levels;
    init_levels(&levels);
    read_yaml_file(argv[1], &levels);

    int grid_rows = 0;
    int grid_cols = 0;

    for (size_t i = 0; i < levels.count; i++) {
        if (grid_rows < levels.items[i].count) {
            grid_rows = levels.items[i].count;
        }

        for (size_t j = 0; j < levels.items[i].count; j++) {
            if (grid_cols < levels.items[i].value_lengths[j]) {
                grid_cols = levels.items[i].value_lengths[j];
            }
        }
    }

    printf("Welcome to Cardinal Chains!\n\n");
    CardinalChainsGame game = init_game(&levels.items[0], 0);

    while (game.current_level < levels.count) {
        play_game(&game, &levels, grid_rows, grid_cols);

        if (game.is_game_over) {
            printf("\nYou didn't complete the level. Try again!\n");
            break;
        } else if (game.current_level < levels.count) {
            printf("\nCongratulations! You've completed the level. Moving on to the "
                   "next level.\n");
        } else {
            printf("Thanks for playing!\n");
        }
    }
    free_levels(&levels);
    free_game(&game);

    return 0;
}
