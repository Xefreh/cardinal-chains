#pragma once

#include <stdbool.h>

#include "level.h"

typedef struct {
    int x;
    int y;
} Position;

typedef enum {
    NORTH, SOUTH, EAST, WEST
} Direction;

typedef struct {
    Level level;
    Position **chains;
    int *chain_lengths;
    int chain_count;
    bool is_game_over;
    int current_level;
} CardinalChainsGame;

int count_chains(Level *level);

CardinalChainsGame init_game(Level *level, int current_level);
void free_game(CardinalChainsGame *game);

bool move_chain(CardinalChainsGame *game, int chain_index, Direction direction,
                int grid_rows, int grid_cols);
void cancel_last_move(CardinalChainsGame *game, int chain_index);
void erase_chain(CardinalChainsGame *game, int chain_index);
void restart_level(CardinalChainsGame *game);

bool is_game_completed(CardinalChainsGame *game);
void load_next_level(CardinalChainsGame *game, Levels *levels);
