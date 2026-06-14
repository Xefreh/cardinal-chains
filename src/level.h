#pragma once

#include <stddef.h>

typedef struct {
    int id;
    int **values;
    size_t *value_lengths;
    size_t count;
} Level;

typedef struct {
    Level *items;
    size_t count;
} Levels;

void init_levels(Levels *levels);
void add_level(Levels *levels, int number, int **values, size_t *value_lengths,
               size_t count);
void free_levels(Levels *levels);
