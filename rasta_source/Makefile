SHELL=/bin/sh
.SUFFIXES:
.SUFFIXES: .cpp .h .o

CXX=g++
CXXFLAGS= -Wall -march=native -std=c++11 
FASTFLAGS=-O3
TITLE=demo

.PHONY : all clean

# make all
all: test mtxtest

test: Rasta.cpp test.cpp *.h
	$(CXX) $(CXXFLAGS) $^ -o $@
mtxtest: Rasta.cpp mtxtest.cpp *.h
	$(CXX) $(CXXFLAGS) $^ -o $@

run: test
	./test

clean:
	-rm -rf test
