
# Online pacman game

## Running game

### Makefile
<strong>Build</strong>
```
make build
```
<strong>Run server<strong>
```
./bin/main -s -p 4444 -n 4 -name "player1"
```
<strong>Run client<strong>
```
./bin/main -ip 127.0.0.1 -p 4444 -name "player2"
```
### Docker
<strong>Build</strong>
```
docker build . -t "online-pacman"
```
<strong>Run server<strong>
```
docker run --net=host -it --rm "online-pacman" /app/bin/main -s -p 4444 -n 2 -name "player1"
```
<strong>Run client<strong>
```
docker run --net=host -it --rm "online-pacman" /app/bin/main -ip 127.0.0.1 -p 4444 -name "player2"
```