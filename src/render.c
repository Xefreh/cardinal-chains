#include "render.h"

#include <stdio.h>

void print_game_grid(const CardinalChainsGame *game, int current_chain,
                     int active_chains[]) {
    const char *chain_colors[] = {"\033[31m", "\033[32m", "\033[33m",
                                  "\033[34m", "\033[35m", "\033[36m"};
    const char *reset_color = "\033[0m";

    for (int i = 0; i < game->level.count; i++) {
        for (int j = 0; j < game->level.value_lengths[i]; j++) {
            bool is_part_of_chain = false;
            int chain_index;

            for (int k = 0; k < game->chain_count; k++) {
                for (int l = 0; l < game->chain_lengths[k]; l++) {
                    if (game->chains[k][l].x == i && game->chains[k][l].y == j) {
                        is_part_of_chain = true;
                        chain_index = active_chains[k];
                        break;
                    }
                }
                if (is_part_of_chain)
                    break;
            }

            if (is_part_of_chain) {
                printf("%s", chain_colors[chain_index % (sizeof(chain_colors) /
                                                         sizeof(chain_colors[0]))]);
                int value = game->level.values[i][j];
                printf("%c ", value == -1 ? 'x' : '0' + value);
                printf("%s", reset_color);
            } else if (game->level.values[i][j] != 0) {
                printf("%d ", game->level.values[i][j]);
            } else {
                printf("  ");
            }
        }
        printf("\n");
    }

    printf("Current position: row %d, column %d\n",
           game->chains[current_chain][game->chain_lengths[current_chain] - 1].x +
           1,
           game->chains[current_chain][game->chain_lengths[current_chain] - 1].y +
           1);
    printf("Current chain color: %sChain %d%s\n",
           chain_colors[active_chains[current_chain] %
                        (sizeof(chain_colors) / sizeof(chain_colors[0]))],
           active_chains[current_chain] + 1, reset_color);
}
