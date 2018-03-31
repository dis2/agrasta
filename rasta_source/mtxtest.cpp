#include "Rasta.h"
#include <iostream>
#include <random>
#include <chrono>
#include <vector>
#include <array>
#include <set>

int main() {
	Rasta r(0xFFD54AA, lu);
	r.initrand(1,2);
	int faults = 0;
	for (int i = 0; i < 1000; i++) {
		std::vector<block> mat;
		mat.clear();
		for (int i = 0; i < blocksize; ++i) {
			mat.push_back(r.getrandblock());
		}

		int k = r.rank_of_Matrix(mat);
		if (k != blocksize) {
			faults++;
		}
	}
	std::cout << faults << "\n";
}

