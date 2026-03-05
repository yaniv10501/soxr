#!/bin/sh
set -e

# Configuration
ARM_OMP=$(brew --prefix libomp)
INTEL_OMP="/usr/local/opt/libomp"

case "$1" in -j*) j="$1"; shift;; esac
test x"$1" = x && build=Release || build="$1"

# Clean previous build artifacts
rm -rf build_arm64 build_x86_64 universal CMakeCache.txt

# 1. Build ARM64 (Native M3)
echo "Building for arm64..."
mkdir -p build_arm64 && cd build_arm64
cmake --fresh -Wno-dev \
  -DCMAKE_BUILD_TYPE="$build" \
  -DCMAKE_OSX_ARCHITECTURES="arm64" \
  -DBUILD_SHARED_LIBS=ON \
  -DCMAKE_CXX_STANDARD=17 \
  -DOpenMP_C_FLAGS="-Xpreprocessor -fopenmp -I$ARM_OMP/include" \
  -DOpenMP_C_LIB_NAMES="omp" \
  -DOpenMP_CXX_FLAGS="-Xpreprocessor -fopenmp -I$ARM_OMP/include" \
  -DOpenMP_CXX_LIB_NAMES="omp" \
  -DOpenMP_omp_LIBRARY="$ARM_OMP/lib/libomp.dylib" \
  -DCMAKE_C_FLAGS="-I$ARM_OMP/include" \
  -DCMAKE_CXX_FLAGS="-I$ARM_OMP/include" \
  -DCMAKE_SHARED_LINKER_FLAGS="-L$ARM_OMP/lib -lomp" \
  -DCMAKE_EXE_LINKER_FLAGS="-L$ARM_OMP/lib -lomp" \
  ..
make $j
cd ..

# 2. Build x86_64 (Intel)
echo "Building for x86_64..."
mkdir -p build_x86_64 && cd build_x86_64
cmake --fresh -Wno-dev \
  -DCMAKE_BUILD_TYPE="$build" \
  -DCMAKE_OSX_ARCHITECTURES="x86_64" \
  -DBUILD_SHARED_LIBS=ON \
  -DCMAKE_CXX_STANDARD=17 \
  -DOpenMP_C_FLAGS="-Xpreprocessor -fopenmp -I$INTEL_OMP/include" \
  -DOpenMP_C_LIB_NAMES="omp" \
  -DOpenMP_CXX_FLAGS="-Xpreprocessor -fopenmp -I$INTEL_OMP/include" \
  -DOpenMP_CXX_LIB_NAMES="omp" \
  -DOpenMP_omp_LIBRARY="$INTEL_OMP/lib/libomp.dylib" \
  -DCMAKE_C_FLAGS="-I$INTEL_OMP/include" \
  -DCMAKE_CXX_FLAGS="-I$INTEL_OMP/include" \
  -DCMAKE_SHARED_LINKER_FLAGS="-L$INTEL_OMP/lib -lomp" \
  -DCMAKE_EXE_LINKER_FLAGS="-L$INTEL_OMP/lib -lomp" \
  ..
make $j
cd ..

# 3. Create Universal Binary
echo "Creating universal binary..."
mkdir -p universal
lipo -create build_arm64/src/libsoxr.dylib build_x86_64/src/libsoxr.dylib -output universal/libsoxr.dylib

echo "Done! Universal binary is in ./universal/"
file universal/libsoxr.dylib