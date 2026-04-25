# filequeue

File system based queue for golang

## Commands

### push

Enqueue a new item to the queue.

```bash
filequeue push -k <kind> -m "<message>"
filequeue push -k exec -m "echo hello"
echo "echo world" | filequeue push -k exec
```

Flags:
- `-d, --queuedir`: Queue directory (default: `~/filequeue`)
- `-k, --kind`: Queue kind - `exec` or `clipboard` (default: `exec`)
- `-m, --message`: Message to queue (read from stdin if omitted)

### pop

Dequeue and execute a single item from the queue.

```bash
filequeue pop
filequeue pop -d /custom/queue/dir
```

Flags:
- `-d, --queuedir`: Queue directory (default: `~/filequeue`)

### consume

Continuously consume and process queue items.

```bash
filequeue consume
filequeue consume -d /custom/queue/dir
filequeue consume -p /custom/path/filequeue.pid -i 10s
```

Flags:
- `-d, --queuedir`: Queue directory (default: `~/filequeue`)
- `-p, --pidfile`: PID file path (default: `~/.filequeue/filequeue.pid`)
- `-i, --interval`: Polling interval (default: `5s`, used as fallback)

The consume command watches the queue directory for new items and automatically processes them. It prevents multiple consumer instances using a PID file. Send SIGTERM or SIGINT to gracefully shutdown.

### purge

Remove all items from the queue.

```bash
filequeue purge
filequeue purge -d /custom/queue/dir
```

Flags:
- `-d, --queuedir`: Queue directory (default: `~/filequeue`)

### version

Show version information.

```bash
filequeue version
```

## Queue Kinds

### exec

Execute a shell command. The message is treated as a command line.

```bash
filequeue push -k exec -m "echo hello world"
filequeue push -k exec -m "sleep 5 && echo done"
```

Environment variables are expanded before execution.

### clipboard

Copy the message to the clipboard.

```bash
filequeue push -k clipboard -m "text to copy"
```

On macOS, uses `pbcopy`. On Linux, uses `cat` (which writes to stdout).

## Examples

### Running as a systemd service

Create a systemd unit file at `~/.config/systemd/user/filequeue.service`:

```ini
[Unit]
Description=filequeue consumer
After=network.target

[Service]
Type=simple
ExecStart=/path/to/filequeue consume -d ~/filequeue -p ~/.filequeue/filequeue.pid
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
```

Then:

```bash
systemctl --user enable filequeue
systemctl --user start filequeue
systemctl --user status filequeue
```

### Using with scripts

Queue commands for later execution:

```bash
#!/bin/bash
while read -r cmd; do
  filequeue push -k exec -m "$cmd"
done

# Start consumer in background
filequeue consume &
```
