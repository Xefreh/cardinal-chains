#include "level_loader.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <yaml.h>

static int *parse_int_sequence(const char *str, size_t *count) {
    int *values = NULL;
    *count = 0;
    char *tmp = strdup(str);
    char *token = strtok(tmp, " ");
    while (token != NULL) {
        (*count)++;
        values = realloc(values, (*count) * sizeof(int));
        values[(*count) - 1] = atoi(token);
        token = strtok(NULL, " ");
    }
    free(tmp);
    return values;
}

void read_yaml_file(const char *filename, Levels *levels) {
    FILE *file = fopen(filename, "r");
    if (!file) {
        fprintf(stderr, "Failed to open file: %s\n", filename);
        return;
    }

    yaml_parser_t parser;
    yaml_event_t event;

    if (!yaml_parser_initialize(&parser)) {
        fprintf(stderr, "Failed to initialize the YAML parser\n");
        fclose(file);
        return;
    }

    yaml_parser_set_input_file(&parser, file);

    int done = 0;
    int level = 0;
    int current_level_number = -1;
    int previous_level_number;
    int **values = NULL;
    size_t *value_lengths = NULL;
    size_t value_count = 0;

    while (!done) {
        if (!yaml_parser_parse(&parser, &event)) {
            fprintf(stderr, "Parser error %d\n", parser.error);
            break;
        }

        if (event.type == YAML_MAPPING_START_EVENT ||
            event.type == YAML_SEQUENCE_START_EVENT) {
            level++;
        } else if (event.type == YAML_MAPPING_END_EVENT ||
                   event.type == YAML_SEQUENCE_END_EVENT) {
            level--;
            if (level == 1 && value_count > 0) {
                add_level(levels, current_level_number, values, value_lengths,
                          value_count);
                values = NULL;
                value_lengths = NULL;
                value_count = 0;
            }
        }

        if (event.type == YAML_SCALAR_EVENT) {
            if (level == 2) {
                previous_level_number = current_level_number;
                current_level_number =
                        atoi((const char *) event.data.scalar.value);

                if (previous_level_number != -1 && value_count > 0) {
                    add_level(levels, previous_level_number, values,
                              value_lengths, value_count);
                    values = NULL;
                    value_lengths = NULL;
                    value_count = 0;
                }
            } else if (level == 3) {
                size_t sequence_length;
                int *sequence = parse_int_sequence(
                        (const char *) event.data.scalar.value,
                        &sequence_length);
                value_count++;
                values = realloc(values, value_count * sizeof(int *));
                values[value_count - 1] = sequence;
                value_lengths =
                        realloc(value_lengths, value_count * sizeof(size_t));
                value_lengths[value_count - 1] = sequence_length;
            }
        }

        done = (event.type == YAML_STREAM_END_EVENT);
        yaml_event_delete(&event);
    }

    if (value_count > 0) {
        add_level(levels, current_level_number, values, value_lengths,
                  value_count);
    }

    yaml_parser_delete(&parser);
    fclose(file);
}
