#include "game_loop.h"

#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>

#include "input.h"
#include "render.h"

void play_game(CardinalChainsGame *game, Levels *levels, int grid_rows,
               int grid_cols) {
    int current_chain = 0;
    char input;

    int active_chains[game->chain_count];
    for (int i = 0; i < game->chain_count; i++) {
        active_chains[i] = i;
    }
    printf("\n");
    printf("Chain count: %d\n", game->chain_count);

    while (true) {
        print_game_grid(game, current_chain, active_chains);
        input = get_user_input();
        getchar();
        printf("\n");

        if (toupper(input) == 'Q') {
            printf("Thank you for playing, see you next time!\n");
            exit(0);
        }

        Direction direction;
        bool moved = false;

        switch (toupper(input)) {
            case 'N':
                direction = NORTH;
                moved = move_chain(game, current_chain, direction, grid_rows,
                                   grid_cols);
                break;
            case 'S':
                direction = SOUTH;
                moved = move_chain(game, current_chain, direction, grid_rows,
                                   grid_cols);
                break;
            case 'E':
                direction = EAST;
                moved = move_chain(game, current_chain, direction, grid_rows,
                                   grid_cols);
                break;
            case 'W':
                direction = WEST;
                moved = move_chain(game, current_chain, direction, grid_rows,
                                   grid_cols);
                break;
            case 'B':
                cancel_last_move(game, current_chain);
                break;
            case 'R':
                erase_chain(game, current_chain);
                break;
            case 'X':
                restart_level(game);
                break;
            case 'C':
                current_chain = (current_chain + 1) % game->chain_count;
                break;
            default:
                printf("Invalid input. Please try again.\n\n");
                continue;
        }

        if (moved && is_game_completed(game)) {
            current_chain = 0;
            print_game_grid(game, current_chain, active_chains);
            printf("\nLevel completed!\n\n");
            game->current_level++;
            if (game->current_level < levels->count) {
                load_next_level(game, levels);
            } else {
                printf("Congratulations! You have completed all levels!\n");
                break;
            }
        }
    }
}
