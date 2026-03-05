@echo off
setlocal enabledelayedexpansion

:: Set build type (Default to Release)
set BUILD_TYPE=Release
if not "%~1"=="" set BUILD_TYPE=%~1

:: Clean previous build
del/f CMakeCache.txt
if exist build_win_x64 rmdir /s /q build_win_x64
mkdir build_win_x64
cd build_win_x64

echo Building for Windows x64 (%BUILD_TYPE%)...

:: Configure with OpenMP enabled (MSVC uses /openmp)
cmake -A x64 ^
      -DCMAKE_BUILD_TYPE=%BUILD_TYPE% ^
      -DBUILD_SHARED_LIBS=ON ^
      -DOpenMP_C_FLAGS="/openmp" ^
      -DOpenMP_CXX_FLAGS="/openmp" ^
      -DCMAKE_C_FLAGS="/W3 /O2" ^
      -DCMAKE_CXX_FLAGS="/W3 /O2" ^
      ..

:: Build the project
cmake --build . --config %BUILD_TYPE%

echo.
echo Build complete.
echo DLL located in: build_win_x64\src\%BUILD_TYPE%\libsoxr.dll
cd ..