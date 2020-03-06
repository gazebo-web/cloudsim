#!/usr/bin/env ruby
require 'socket'
require 'yaml'

class String
  def is_int?
    true if Integer(self) rescue false
  end
end

# A class that wraps the tc linux command line tool. Usage example:
#
#   # Create and run the traffic shaper
#   TrafficShaper.new('my_yaml_file.yml')
#
#   # YAML file format. Numerical values can be either a single value, or
#   a 2-value array. A 2-value array will be used as a range from which a 
#   random number is chosen.
#
#   - <name_for_a_shaper_instance>:
#     # Add a delay rule.
#     delay:
#       # Delay duration in milliseconds.
#       duration_ms: <number, or array with 2-values> 
#
#       # Jitter duration in milliseconds.
#       jitter_ms: <number, or array with 2-values> 
#
#       # Correlation percentage.
#       correlation_percent: <number, or array with 2-values> 
#
#     # Add a packet rule.
#     packet:
#       # Packet loss percentage.
#       loss_percent: <number, or array with 2-values>
#
#       # Packet corruption percentage.
#       corrupt_percent: <number, or array with 2-values>
#
#       # Packet duplication percentage.
#       duplicate_percent: <number, or array with 2-values>
#
#     # Add a bandwidth rule.
#     bandwidth:
#       # Maximum sustained bit rate in kilobytes per second.
#       max_rate_kbps: <number, or array with 2-values>
#
#       # Maximum allowed burst in kilobytes.
#       burst_kb: <number, or array with 2-values>
#
#       # Milliseconds a packet can be in the queue before it gets dropped.
#       latency_ms: <number, or array with 2-values>
#
#     # Duration of a run in seconds.
#     duration_s: <number, or array with 2-values>
#
#     # Number of iterations to run this shaper instance.
#     iterations: <number, or array with 2-values>

class TrafficShaper

  # Initialize the traffic shaper.
  def initialize(file)
    @delay_cmd = ""
    @loss_cmd = ""
    @bandwidth_cmd = ""

    printf "# Interfaces\n"
    # Get all the interfaces.
    @ifaces = []
    addrs = Socket.getifaddrs
    addrs.each do |addr_info|
      if addr_info.addr && addr_info.addr.ipv4?
        @ifaces.push(addr_info.name)
        puts "  * #{addr_info.name}"
      end
    end

    shape_data = YAML.load(File.read(file))
    shape_data.each do |data|
      reset()
      iterations = [1, parse_value(data, "iterations")].max
      puts "\n# Traffic Shaper #{data.keys[0]}"
      for i in 1..iterations
        puts "  * Iteration #{i}"
        if data.has_key?("delay")
          duration_ms = parse_value(data["delay"], "duration_ms")
          jitter_ms = parse_value(data["delay"], "jitter_ms")
          correlation_percent = parse_value(data["delay"], "correlation_percent")
          add_delay(duration_ms, jitter_ms, correlation_percent)
        end
    
        if data.has_key?("packet")
          loss_percent = parse_value(data["packet"], "loss_percent")
          corrupt_percent = parse_value(data["packet"], "corrupt_percent")
          duplication_percent = parse_value(data["packet"], "duplicate_percent")
          add_packet(loss_percent, corrupt_percent, duplication_percent)
        end
    
        if data.has_key?("bandwidth")
          max_rate_kbps = parse_value(data["bandwidth"], "max_rate_kbps")
          burst_kb = parse_value(data["bandwidth"], "burst_kb")
          latency_ms = parse_value(data["bandwidth"], "latency_ms")
          add_bandwidth(max_rate_kbps, burst_kb, latency_ms)
        end
    
        duration_s = parse_value(data, "duration_s")
        run(duration_s)
      end
    end
  end

  # Add a network delay for a specific duration.
  # \param[in] delay Delay duration in milliseconds.
  # \param[in] delay Jitter duration in milliseconds.
  # \param[in] correlation Correlation percentage, between 0 and 100.
  def add_delay(delay_ms, jitter_ms, correlation)
    if delay_ms > 0
      @delay_cmd = " delay #{delay_ms}ms "
      if jitter_ms > 0
        @delay_cmd += " #{jitter_ms}ms "
        # Correlation requires jitter.
        if correlation > 0
          @delay_cmd += " #{[100, correlation].min}% " 
        end
      end
      @delay_cmd += " distribution normal "
    end
  end

  # Add packet loss, corruption and duplication.
  # \param[in] loss Packet loss percentage.
  # \param[in] corrupt Packet corruption percentage.
  # \param[in] duplicate Packet duplication percentage.
  def add_packet(loss, corrupt, duplicate)
    if loss > 0
      @loss_cmd += " loss #{[loss, 100].min}% "
    end
    if corrupt > 0
      @loss_cmd += " corrupt #{[corrupt, 100].min}% "
    end
    if duplicate > 0
      @loss_cmd += " duplicate #{[duplicate, 100].min}% "
    end
  end

  # Add bandwidth limit
  # \param[in] max_rate_kbps Maximum sustained bit rate in kilobytes per second.
  # \param[in] burst_kbps Maximum allowed burst in kilobytes.
  # \param[in] latency_ms Milliseconds a packet can be in the queue before
  # it gets dropped.
  def add_bandwidth(max_rate_kbps, burst_kb, latency_ms)
    if max_rate_kbps > 0 && burst_kb > 0 && latency_ms > 0
      @bandwidth_cmd = " tbf rate #{max_rate_kbps}kbps burst #{burst_kb}k " +
        "latency #{latency_ms}ms"
    end
  end

  # Run the set of rules on all interfaces for a duration. This function
  # blocks for the specified duration.
  # \param[in] duration Number of seconds to run the rules.
  def run(duration)
    # Clear all the rules
    clear_rules()
    begin
      # Output some useful information.
      printf "\t# Duration %ds\n", duration
      printf "\t  * Delay[%s]\n", @delay_cmd
      printf "\t  * Loss[%s]\n", @loss_cmd
      printf "\t  * Bandwidth[%s]\n", @bandwidth_cmd

      @ifaces.each do |iface|
        if @delay_cmd != "" || @loss_cmd != ""
          system("tc qdisc add dev #{iface} root handle 1:0 " +
                 "netem #{@delay_cmd} #{@loss_cmd}")
        end

        if @bandwidth_cmd != ""
          system("tc qdisc add dev #{iface} parent 1:1 handle 10:0 " +
                 "#{@bandwidth_cmd}")
        end
      end
      sleep duration
    ensure
      reset()
    end
  end

  # Clear all the parameters and rules.
  def reset()
    @delay_cmd = ""
    @loss_cmd = ""
    @bandwidth_cmd = ""

    clear_rules()
  end

  # Clear just the rules sent to tc.
  def clear_rules()
    # Clear all the rules.
    @ifaces.each do |iface|
      system("tc qdisc del dev #{iface} root 2>/dev/null")
    end
  end

  # Parse a value from a hash generated from a yaml file.
  def parse_value(hash, key)
    result = hash.has_key?(key) ? hash[key] : 0
  
    # Convert a string to a number
    if result.is_a?(String) && result.is_int?
      result = result.to_i
      # Convert an array to a random number
    elsif result.is_a? Array
      if result.length == 2
        result = rand(result[0]..result[1])
      else
        printf "Array must have a length of 2."
        result = 0
      end
    end
  
    return result
  end
end

# Create and run the traffic shaper
TrafficShaper.new("shaper.yml")
