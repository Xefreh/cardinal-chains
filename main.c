#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <yaml.h>
#include <stdbool.h>

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

typedef struct {
	int x;
	int y;
} Position;

typedef struct {
	Level level;
	Position **chains;
	int *chain_lengths;
	int chain_count;
	bool is_game_over;
	int current_level;
} CardinalChainsGame;

typedef enum {
	NORTH,
	SOUTH,
	EAST,
	WEST
} Direction;

void init_levels(Levels *levels) {
	levels->items = NULL;
	levels->count = 0;
}

void add_level(Levels *levels, int number, int **values, size_t *value_lengths, size_t count) {
	levels->count++;
	levels->items = realloc(levels->items, levels->count * sizeof(Level));

	Level *level = &levels->items[levels->count - 1];
	level->id = number;
	level->values = values;
	level->value_lengths = value_lengths;
	level->count = count;

//	printf("Added level %d:\n", level->id);
//	for (size_t i = 0; i < level->count; i++) {
//		printf("Sequence: ");
//		for (size_t j = 0; j < level->value_lengths[i]; j++) {
//			printf("%d ", level->values[i][j]);
//		}
//		printf("\n");
//	}
//	printf("\n");
}

int *parse_int_sequence(const char *str, size_t *count) {
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

		if (event.type == YAML_MAPPING_START_EVENT || event.type == YAML_SEQUENCE_START_EVENT) {
			level++;
		} else if (event.type == YAML_MAPPING_END_EVENT || event.type == YAML_SEQUENCE_END_EVENT) {
			level--;
			if (level == 1 && value_count > 0) {
				add_level(levels, current_level_number, values, value_lengths, value_count);
				values = NULL;
				value_lengths = NULL;
				value_count = 0;
			}
		}

		if (event.type == YAML_SCALAR_EVENT) {
			if (level == 2) {
				previous_level_number = current_level_number;
				current_level_number = atoi((const char *) event.data.scalar.value);

				if (previous_level_number != -1 && value_count > 0) {
					add_level(levels, previous_level_number, values, value_lengths, value_count);
					values = NULL;
					value_lengths = NULL;
					value_count = 0;
				}

//				printf("Parsed level id: %d\n", current_level_number);
			} else if (level == 3) {
				size_t sequence_length;
				int *sequence = parse_int_sequence((const char *) event.data.scalar.value, &sequence_length);
				value_count++;
				values = realloc(values, value_count * sizeof(int *));
				values[value_count - 1] = sequence;
				value_lengths = realloc(value_lengths, value_count * sizeof(size_t));
				value_lengths[value_count - 1] = sequence_length;

//				printf("Parsed sequence for level %d: ", current_level_number);
//				for (size_t i = 0; i < sequence_length; i++) {
//					printf("%d ", sequence[i]);
//				}
//				printf("\n");
			}
		}

		done = (event.type == YAML_STREAM_END_EVENT);
		yaml_event_delete(&event);
	}

	if (value_count > 0) {
		add_level(levels, current_level_number, values, value_lengths, value_count);
	}

	yaml_parser_delete(&parser);
	fclose(file);
}

//void print_levels(const Levels *levels) {
//	for (size_t i = 0; i < levels->count; i++) {
//		Level *level = &levels->items[i];
//		printf("Level %d:\n", level->id);
//		for (size_t j = 0; j < level->count; j++) {
//			int *sequence = level->values[j];
//			size_t sequence_length = level->value_lengths[j];
//			for (size_t k = 0; k < sequence_length; k++) {
//				printf("%d ", sequence[k]);
//			}
//			printf("\n");
//		}
//		printf("\n");
//	}
//}

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

int find_next_empty_chain(const CardinalChainsGame *game) {
	for (int i = 0; i < game->chain_count; i++) {
		if (game->chain_lengths[i] == 0) {
			return i;
		}
	}

	return -1;
}

void init_chains(CardinalChainsGame *game) {
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

char get_user_input() {
	printf("Select a direction (N, S, E, W).\n");
	printf("Cancel the previous move (B).\n");
	printf("Erase the chain (R).\n");
	printf("Restart the level (X).\n");
	printf("Select another chain (C).\n");
	return getchar();
}

void print_game_grid(const CardinalChainsGame *game, int current_chain, int active_chains[]) {
	const char *chain_colors[] = {"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m"};
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

//						printf("Chain index: %d, active chain: %d\n", k, chain_index);
						break;
					}
				}
				if (is_part_of_chain) break;
			}

			if (is_part_of_chain) {
				printf("%s", chain_colors[chain_index % (sizeof(chain_colors) / sizeof(chain_colors[0]))]);
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
		   game->chains[current_chain][game->chain_lengths[current_chain] - 1].x + 1,
		   game->chains[current_chain][game->chain_lengths[current_chain] - 1].y + 1);
	printf("Current chain color: %sChain %d%s\n",
		   chain_colors[active_chains[current_chain] % (sizeof(chain_colors) / sizeof(chain_colors[0]))],
		   active_chains[current_chain] + 1, reset_color);
}

bool move_chain(CardinalChainsGame *game, int chain_index, Direction direction, int grid_rows, int grid_cols) {
	if (chain_index < 0 || chain_index >= game->chain_count) {
		return false;
	}

	Position last_position = game->chains[chain_index][game->chain_lengths[chain_index] - 1];
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
			if (game->chains[i][j].x == new_position.x && game->chains[i][j].y == new_position.y) {
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
	game->chains[chain_index] = realloc(game->chains[chain_index], game->chain_lengths[chain_index] * sizeof(Position));
	game->chains[chain_index][game->chain_lengths[chain_index] - 1] = new_position;

	return true;
}

void cancel_last_move(CardinalChainsGame *game, int chain_index) {
	int chain_length = game->chain_lengths[chain_index];

	if (chain_length > 1) {
		game->chain_lengths[chain_index]--;

		game->chains[chain_index] = (Position *) realloc(game->chains[chain_index],
														 (chain_length - 1) * sizeof(Position));
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

void play_game(CardinalChainsGame *game, Levels *levels, int grid_rows, int grid_cols) {
	int current_chain = 0;
	char input;

	int active_chains[game->chain_count];
//	printf("Active chains: ");
	for (int i = 0; i < game->chain_count; i++) {
//		printf("%d ", active_chains[i]);
		active_chains[i] = i;
	}
	printf("\n");
	printf("Chain count: %d\n", game->chain_count);

//	printf("Before move_chain: game->chains[current_chain][game->chain_lengths[current_chain] - 1].x: %d\n",
//		   game->chains[current_chain][game->chain_lengths[current_chain] - 1].x + 1);
//	printf("Before move_chain: game->chains[current_chain][game->chain_lengths[current_chain] - 1].y: %d\n",
//		   game->chains[current_chain][game->chain_lengths[current_chain] - 1].y + 1);

	while (true) {
		print_game_grid(game, current_chain, active_chains);
		input = get_user_input();
		getchar();
		printf("\n");

		if (input == 'Q') {
			printf("Thank you for playing, see you next time!\n");
			exit(0);
		}

		Direction direction;
		bool moved = false;

		switch (input) {
			case 'N':
				direction = NORTH;
				moved = move_chain(game, current_chain, direction, grid_rows, grid_cols);
				break;
			case 'S':
				direction = SOUTH;
				moved = move_chain(game, current_chain, direction, grid_rows, grid_cols);
				break;
			case 'E':
				direction = EAST;
				moved = move_chain(game, current_chain, direction, grid_rows, grid_cols);
				break;
			case 'W':
				direction = WEST;
				moved = move_chain(game, current_chain, direction, grid_rows, grid_cols);
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

//		if (moved) {
//			printf("After move_chain: game->chains[current_chain][game->chain_lengths[current_chain] - 1].x: %d\n",
//				   game->chains[current_chain][game->chain_lengths[current_chain] - 1].x);
//			printf("After move_chain: game->chains[current_chain][game->chain_lengths[current_chain] - 1].y: %d\n",
//				   game->chains[current_chain][game->chain_lengths[current_chain] - 1].y);
//		}

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

void free_cardinal_chains_game(CardinalChainsGame *game) {
	for (int i = 0; i < game->chain_count; i++) {
		free(game->chains[i]);
	}
	free(game->chains);
	free(game->chain_lengths);
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

int main(int argc, char *argv[]) {
	if (argc < 2) {
		fprintf(stderr, "Usage: %s <YAML file>\n", argv[0]);
		return 1;
	}

	Levels levels;
	init_levels(&levels);
	read_yaml_file(argv[1], &levels);

//	print_levels(&levels);

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
			printf("\nCongratulations! You've completed the level. Moving on to the next level.\n");
		} else {
			printf("Thanks for playing!\n");
		}
	}
	free_levels(&levels);
	free_cardinal_chains_game(&game);

	return 0;
}