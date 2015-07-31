#include <fstream>
#include <iostream>

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

    if (letter1 > 0) std::cout << static_cast<char>(64 + letter1);
    if (letter2 > 0) std::cout << static_cast<char>(64 + letter2);
    std::cout << static_cast<char>(64 + letter3) << std::endl;
}

int main(int argc, char* argv[]) {
    std::ifstream inputFile(argv[1]);

    int columnNumber;
    while (inputFile >> columnNumber) {
        printExcelColumnName(columnNumber);
    }
}
