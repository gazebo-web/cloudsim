# Network traffic shaping

## Local usage

1. Install `iproute2` and `ruby2.5`. On Ubuntu:

```
sudo apt-get install iproute2 ruby2.5
```

1. In a terminal, ping something.

```
ping google.com
```

1. Run the shaper in another terminal, and watch the ping output.

```
sudo ./shaper.rb
```

## Docker usage

1. Install docker.

1. Build the `shaper` container. 

```
./build.bash
```

1. In a terminal, ping something.

```
ping google.com
```

2. Run the docker container, and watch the ping output.

```
./run.bash
```

## YAML configuration

The `shaper.rb` script reads the `shaper.yml` file, and executes the `tc`
command based on its contents. Refer to the documentation in `shaper.rb` for
information about the YAML format used by `shaper`.
