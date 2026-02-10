git clone --recurse-submodules https://github.com/Rocco2300/particle-simulation.git

cd particle-simulation

mkdir build

cd build

cmake .. -G "MinGW Makefiles"

cmake --build .

cd ../..

copy "particle-simulation\build\libraylib.dll" "build\libraylib.dll"
copy "particle-simulation\build\libphysics.dll" "build\libphysics.dll"

mkdir build

go env -w GOBIN=%~dp0build

cd client/game
go install

cd ../console
go install

cd ../../server
go install
