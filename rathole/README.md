# Rathole VPN Tunnel

Rathole is a fast, stable, and easy-to-use reverse proxy tunnel that can expose a local service to the internet through a server.

## Overview

This directory contains:
- `rathole` - The rathole binary executable
- `v4-client.toml` - IPv4 client configuration
- `v6-client.toml` - IPv6 client configuration
- `sw-v4-server.toml` - IPv4 server configuration
- `sw-v6-server.toml` - IPv6 server configuration

## Configuration

### Before First Use

**IMPORTANT**: Update the `default_token` in all configuration files with a secure token before using.

```bash
# Generate a secure token (optional)
openssl rand -base64 32
```

Then edit each `.toml` file and replace `change_this_token` with your generated token.

### Configuration Files

- **Client configs** (`v4-client.toml`, `v6-client.toml`): Connect to remote server and forward local services
- **Server configs** (`sw-v4-server.toml`, `sw-v6-server.toml`): Listen on remote server and forward to local services

## Usage

### Running as Server

On your remote server, run:

```bash
# IPv4 server
./rathole -s sw-v4-server.toml

# IPv6 server
./rathole -s sw-v6-server.toml
```

### Running as Client

On your local machine, run:

```bash
# IPv4 client
./rathole -c v4-client.toml

# IPv6 client
./rathole -c v6-client.toml
```

### Debugging and Logging

Rathole uses Rust's logging system. You can control log verbosity using the `RUST_LOG` environment variable.

**Log Levels** (from least to most verbose):
- `error` - Only error messages
- `warn` - Warnings and errors
- `info` - Informational messages, warnings, and errors
- `debug` - Debug information, info, warnings, and errors
- `trace` - Very verbose trace information (all logs)

**Examples:**

```bash
# Run with error-level logging only
RUST_LOG=error ./rathole -s sw-v4-server.toml

# Run with debug-level logging
RUST_LOG=debug ./rathole -c v4-client.toml

# Run with trace-level logging (most verbose)
RUST_LOG=trace ./rathole -s sw-v4-server.toml

# Set log level for specific modules
RUST_LOG=rathole=debug,info ./rathole -s sw-v4-server.toml
```

**Default behavior**: If `RUST_LOG` is not set, rathole will use its default logging level (typically `info`).

## Running in Background

### Method 1: Using screen

```bash
# Install screen if not available
# Ubuntu/Debian: sudo apt install screen
# CentOS/RHEL: sudo yum install screen

# Start a new screen session
screen -S rathole-server

# Run rathole (with optional debug logging)
RUST_LOG=debug ./rathole -s sw-v4-server.toml

# Detach: Press Ctrl+A, then D
# Reattach: screen -r rathole-server
# List sessions: screen -ls
```

### Method 2: Using tmux

```bash
# Install tmux if not available
# Ubuntu/Debian: sudo apt install tmux
# CentOS/RHEL: sudo yum install tmux

# Start a new tmux session
tmux new -s rathole-server

# Run rathole (with optional debug logging)
RUST_LOG=debug ./rathole -s sw-v4-server.toml

# Detach: Press Ctrl+B, then D
# Reattach: tmux attach -t rathole-server
# List sessions: tmux ls
```

## Systemd Service (Auto-restart on Boot)

### Create Systemd Service for Server

Create a systemd service file:

```bash
sudo nano /etc/systemd/system/sw-server-v4.service
```

Add the following content (adjust paths as needed):

```ini
[Unit]
Description=Rathole VPN Tunnel Server (SW IPv4)
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/orc-cli/rathole
ExecStart=/root/orc-cli/rathole/rathole -s /root/orc-cli/rathole/sw-v4-server.toml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Note**: To enable debug logging in systemd, change `Environment="RUST_LOG=info"` to `Environment="RUST_LOG=debug"` or `Environment="RUST_LOG=trace"` for more verbose output.

For IPv6 server, create `rathole-server-v6.service` with the same content but change:
- Description to "Rathole VPN Tunnel Server (IPv6)"
- ExecStart to use `sw-v6-server.toml`

### Create Systemd Service for Client

```bash
sudo nano /etc/systemd/system/rathole-client-v4.service
```

```ini
[Unit]
Description=Rathole VPN Tunnel Client (IPv4)
After=network.target

[Service]
Type=simple
User=your_username
WorkingDirectory=/<PATH>/orc-cli/rathole
Environment="RUST_LOG=info"
ExecStart=/<PATH>/orc-cli/rathole/rathole -c /<PATH>/orc-cli/rathole/v4-client.toml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Note**: To enable debug logging in systemd, change `Environment="RUST_LOG=info"` to `Environment="RUST_LOG=debug"` or `Environment="RUST_LOG=trace"` for more verbose output.

### Enable and Start Services

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable rathole-server-v4.service
sudo systemctl enable rathole-client-v4.service

# Start service
sudo systemctl start rathole-server-v4.service
sudo systemctl start rathole-client-v4.service

# Check status
sudo systemctl status rathole-server-v4.service
sudo systemctl status rathole-client-v4.service

# View logs
sudo journalctl -u rathole-server-v4.service -f
sudo journalctl -u rathole-client-v4.service -f
```

## Manual Restart

### If Running with nohup

```bash
# Find and kill the process
pkill -f rathole

# Start again (with optional debug logging)
nohup env RUST_LOG=debug ./rathole -s sw-v4-server.toml > rathole-server-v4.log 2>&1 &
```

### If Running with systemd

```bash
# Restart service
sudo systemctl restart rathole-server-v4.service
sudo systemctl restart rathole-client-v4.service

# Stop service
sudo systemctl stop rathole-server-v4.service

# Start service
sudo systemctl start rathole-server-v4.service
```

### If Running in screen/tmux

```bash
# Reattach to session
screen -r rathole-server
# or
tmux attach -t rathole-server

# Stop with Ctrl+C, then restart
./rathole -s sw-v4-server.toml
```

## Service Management Commands

### Systemd Commands Reference

```bash
# Start
sudo systemctl start rathole-server-v4.service

# Stop
sudo systemctl stop rathole-server-v4.service

# Restart
sudo systemctl restart rathole-server-v4.service

# Reload (if config changed)
sudo systemctl reload rathole-server-v4.service

# Status
sudo systemctl status rathole-server-v4.service

# Enable (start on boot)
sudo systemctl enable rathole-server-v4.service

# Disable (don't start on boot)
sudo systemctl disable rathole-server-v4.service

# View logs (last 50 lines)
sudo journalctl -u rathole-server-v4.service -n 50

# Follow logs (live)
sudo journalctl -u rathole-server-v4.service -f

# View logs since boot
sudo journalctl -u rathole-server-v4.service -b
```

## Troubleshooting

### Check if rathole is running

```bash
ps aux | grep rathole
netstat -tulpn | grep rathole
# or
ss -tulpn | grep rathole
```

### Check logs

```bash
# If using nohup
tail -f rathole-server-v4.log

# If using systemd
sudo journalctl -u rathole-server-v4.service -f

# View logs with filtering (if using debug/trace logging)
sudo journalctl -u rathole-server-v4.service -f | grep -i error
sudo journalctl -u rathole-server-v4.service -f | grep -i debug
```

### Verify configuration

```bash
# Test configuration syntax
./rathole -s sw-v4-server.toml --check-config

# Test with debug logging to see detailed configuration loading
RUST_LOG=debug ./rathole -s sw-v4-server.toml --check-config
```

### Common Issues

1. **Port already in use**: Check if another process is using the port
   ```bash
   lsof -i :4003
   ```

2. **Permission denied**: Make sure rathole binary is executable
   ```bash
   chmod +x rathole
   ```

3. **Connection refused**: Verify firewall settings and that server is accessible
   ```bash
   # Test connection
   telnet myserver.com 4003
   ```

4. **Token mismatch**: Ensure client and server use the same token

### Debugging with RUST_LOG

When troubleshooting issues, enable debug logging to get more detailed information:

```bash
# Enable debug logging for detailed diagnostics
RUST_LOG=debug ./rathole -s sw-v4-server.toml

# Enable trace logging for maximum verbosity (very detailed)
RUST_LOG=trace ./rathole -c v4-client.toml

# Set error-level only to reduce noise
RUST_LOG=error ./rathole -s sw-v4-server.toml
```

**Log level recommendations:**
- **Production**: `RUST_LOG=error` or `RUST_LOG=warn` (minimal logging)
- **Normal operation**: `RUST_LOG=info` (default, balanced)
- **Troubleshooting**: `RUST_LOG=debug` (detailed diagnostics)
- **Deep debugging**: `RUST_LOG=trace` (very verbose, all logs)

**For systemd services**, edit the service file and modify the `Environment` line:
```ini
Environment="RUST_LOG=debug"  # Change from info to debug/trace
```

Then reload and restart:
```bash
sudo systemctl daemon-reload
sudo systemctl restart rathole-server-v4.service
```

## Notes

- Make sure the `default_token` matches between client and server configurations
- Update `remote_addr` in client configs with your actual server address
- Adjust `local_addr` ports according to your needs
- For production use, consider using systemd services for automatic restarts and logging
