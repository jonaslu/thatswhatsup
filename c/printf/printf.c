#include <unistd.h>
#include <stdlib.h>

#include <stdio.h>

const int STARTING_STRING_BUFFER_LEN = 20;

struct string_w_size {
    char* string;
    int size;
};

static void print_string(struct string_w_size *str) {
    write(1, str->string, str->size);
}

static void reverse_string(struct string_w_size *str) {
    int firstIndex = 0;
    int lastIndex = str->size - 1;

    while(firstIndex < lastIndex) {
        char tempValue = str->string[firstIndex];
        str->string[firstIndex++] = str->string[lastIndex];
        str->string[lastIndex--] = tempValue;
    }
}

static void add_char_to_string(char add_value, struct string_w_size *str) {
    if ((str->size % STARTING_STRING_BUFFER_LEN) == 0) {
        str->string = realloc(str->string, str->size + STARTING_STRING_BUFFER_LEN);
    }

    str->string[str->size] = add_value;
    str->size = str->size + 1;
}


static void concatenate_string(struct string_w_size *dest, struct string_w_size *src) {
    int srcIndex = 0;

    while(srcIndex < src->size) {
        add_char_to_string(src->string[srcIndex++], dest);
    }
}

static void print_int_to_string(int value, struct string_w_size *str) {
    char print_index[] = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'};

    while(value > 0) {
        int subtractValue = (value / 10) * 10;

        int printValue = value - subtractValue;
        add_char_to_string(print_index[printValue], str);

        value = value / 10;
    }

    reverse_string(str);
}

static void print_int(int value, struct string_w_size *dest) {
    struct string_w_size int_string = { NULL, 0 };
    print_int_to_string(value, &int_string);
    concatenate_string(dest, &int_string);
}

static void my_printf(char *fmt, int value) {
    struct string_w_size print_buffer = { NULL, 0 };

    int fmt_index = 0;
    char next_token;

    while(next_token = fmt[fmt_index++], next_token != '\0') {
        if (next_token == '%') {
            char formatting_char = fmt[fmt_index++];

            switch (formatting_char) {
                case 'd':
                    print_int(value, &print_buffer);
                break;
                default: exit(1);
            }
        } else {
            add_char_to_string(next_token, &print_buffer);
        }
    }

    print_string(&print_buffer);
}

int main(int argc, char *argv[]) {
    my_printf("abcedf %d abc", 125);
}
