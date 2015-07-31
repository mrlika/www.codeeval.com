#include <stdio.h>
#include <stdlib.h>

void printExcelColumnName(int columnNumber) {
    int letter1 = columnNumber / (26 * 26);
    if (letter1 > 26) {
        letter1 = 26;
    }

    int letter2 = (columnNumber - letter1 * 26 * 26) / 26;
    if (letter2 == 0) {
        if (letter1 != 0) {
            letter2 = 26;
            letter1--;
        }
    } else if (letter2 == 27) {
        letter2 = 26;
    }

    int letter3 = columnNumber - letter1 * 26 * 26 - letter2 * 26;
    if (letter3 == 0) {
        letter3 = 26;
        letter2--;
    }

    if (letter1 > 0) putchar('A' - 1 + letter1);
    if (letter2 > 0) putchar('A' - 1 + letter2);
    putchar('A' - 1 + letter3);
    putchar('\n');
}

int main(int argc, const char* argv[]) {
    FILE *file = fopen(argv[1], "r");
    char line[8];
    while (fgets(line, 1024, file)) {
        printExcelColumnName(atoi(line));
    }
    return 0;
}
