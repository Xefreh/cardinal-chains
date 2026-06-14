#include "game.h"

#include <stdlib.h>

static int find_next_empty_chain(const CardinalChainsGame *game) {
    for (int i = 0; i < game->chain_count; i++) {
        if (game->chain_lengths[i] == 0) {
            return i;
        }
    }

    return -1;
}

static void init_chains(CardinalChainsGame *game) {
    int chain_count = count_chains(&game->level);
    game->chain_count = chain_count;
    game->chains = (Position **) malloc(chain_count * sizeof(Position *));
    game->chain_lengths = (int *) malloc(chain_count * sizeof(int));

    for (int i = 0; i < chain_count; i++) {
        game->chains[i] = NULL;
        game->chain_lengths[i] = 0;
    }

    for (int i = 0; i < game->level.count; i++) {
        for (int j = 0; j < game->level.value_lengths[i]; j++) {
            int value = game->level.values[i][j];
            if (value == -1) {
                int chain_index = find_next_empty_chain(game);
                game->chains[chain_index] = (Position *) malloc(sizeof(Position));
                game->chains[chain_index][0].x = i;
                game->chains[chain_index][0].y = j;
                game->chain_lengths[chain_index] = 1;
            }
        }
    }
}

int count_chains(Level *level) {
    int max_chain = 0;

    for (int i = 0; i < level->count; i++) {
        for (int j = 0; j < level->value_lengths[i]; j++) {
            if (level->values[i][j] == -1) {
                max_chain++;
            }
        }
    }

    return max_chain;
}

CardinalChainsGame init_game(Level *level, int current_level) {
    CardinalChainsGame game;
    game.level = *level;
    game.chain_count = count_chains(level);
    game.chains = malloc(game.chain_count * sizeof(Position *));
    game.chain_lengths = malloc(game.chain_count * sizeof(int));
    game.is_game_over = false;
    game.current_level = current_level;
    init_chains(&game);
    return game;
}

void free_game(CardinalChainsGame *game) {
    for (int i = 0; i < game->chain_count; i++) {
        free(game->chains[i]);
    }
    free(game->chains);
    free(game->chain_lengths);
}

bool move_chain(CardinalChainsGame *game, int chain_index, Direction direction,
                int grid_rows, int grid_cols) {
    if (chain_index < 0 || chain_index >= game->chain_count) {
        return false;
    }

    Position last_position =
            game->chains[chain_index][game->chain_lengths[chain_index] - 1];
    Position new_position = last_position;

    switch (direction) {
        case NORTH:
            if (last_position.x <= 0) {
                return false;
            }
            new_position.x--;
            break;
        case SOUTH:
            if (last_position.x >= grid_rows - 1) {
                return false;
            }
            new_position.x++;
            break;
        case EAST:
            if (last_position.y >= grid_cols - 1) {
                return false;
            }
            new_position.y++;
            break;
        case WEST:
            if (last_position.y <= 0) {
                return false;
            }
            new_position.y--;
            break;
    }

    for (int i = 0; i < game->chain_count; i++) {
        for (int j = 0; j < game->chain_lengths[i]; j++) {
            if (game->chains[i][j].x == new_position.x &&
                game->chains[i][j].y == new_position.y) {
                return false;
            }
        }
    }

    int last_move_value = game->level.values[last_position.x][last_position.y];
    int next_move_value = game->level.values[new_position.x][new_position.y];

    if (next_move_value < last_move_value || next_move_value == 0) {
        return false;
    }

    game->chain_lengths[chain_index]++;
    game->chains[chain_index] =
            realloc(game->chains[chain_index],
                    game->chain_lengths[chain_index] * sizeof(Position));
    game->chains[chain_index][game->chain_lengths[chain_index] - 1] =
            new_position;

    return true;
}

void cancel_last_move(CardinalChainsGame *game, int chain_index) {
    int chain_length = game->chain_lengths[chain_index];

    if (chain_length > 1) {
        game->chain_lengths[chain_index]--;

        game->chains[chain_index] = (Position *) realloc(
                game->chains[chain_index], (chain_length - 1) * sizeof(Position));
    }
}

void erase_chain(CardinalChainsGame *game, int chain_index) {
    while (game->chain_lengths[chain_index] > 1) {
        cancel_last_move(game, chain_index);
    }
}

void restart_level(CardinalChainsGame *game) {
    for (int i = 0; i < game->chain_count; i++) {
        erase_chain(game, i);
    }
}

bool is_game_completed(CardinalChainsGame *game) {
    for (int i = 0; i < game->level.count; i++) {
        for (int j = 0; j < game->level.value_lengths[i]; j++) {
            if (game->level.values[i][j] == -1 || game->level.values[i][j] == 0) {
                continue;
            }

            bool cell_filled = false;

            for (int k = 0; k < game->chain_count; k++) {
                for (int l = 0; l < game->chain_lengths[k]; l++) {
                    if (game->chains[k][l].x == i && game->chains[k][l].y == j) {
                        cell_filled = true;
                        break;
                    }
                }
                if (cell_filled) {
                    break;
                }
            }

            if (!cell_filled) {
                return false;
            }
        }
    }

    return true;
}

void load_next_level(CardinalChainsGame *game, Levels *levels) {
    if (game->current_level >= levels->count) {
        game->is_game_over = true;
        return;
    }
    game->level = levels->items[game->current_level];
    init_chains(game);
}
