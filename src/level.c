#include "level.h"

#include <stdlib.h>

void init_levels(Levels *levels) {
    levels->items = NULL;
    levels->count = 0;
}

void add_level(Levels *levels, int number, int **values, size_t *value_lengths,
               size_t count) {
    levels->count++;
    levels->items = realloc(levels->items, levels->count * sizeof(Level));

    Level *level = &levels->items[levels->count - 1];
    level->id = number;
    level->values = values;
    level->value_lengths = value_lengths;
    level->count = count;
}

void free_levels(Levels *levels) {
    for (size_t i = 0; i < levels->count; i++) {
        for (size_t j = 0; j < levels->items[i].count; j++) {
            free(levels->items[i].values[j]);
        }
        free(levels->items[i].values);
        free(levels->items[i].value_lengths);
    }

    free(levels->items);
}
