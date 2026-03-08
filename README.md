# Operational Transformation (OT) Collaborative Editor

A simple real-time collaborative text editor demonstrating Operational Transformation with WebSockets, built in Go.

## What is Operational Transformation?

Operational Transformation (OT) is a technology for supporting real-time collaborative editing. It allows multiple users to edit the same document simultaneously while maintaining consistency across all copies.

### The Core Problem

When two users edit the same document at the same time, their operations can conflict:

```
Initial document: "Hello"
User A at position 5: Insert "!"  → "Hello!"
User B at position 0: Insert "Hey " → "Hey Hello"
```

If we just apply these operations in order, we get inconsistent results depending on the order.

### How OT Solves It

OT uses **transformation functions** to adjust operations based on concurrent operations:

1. Each operation has a position and an action (insert/delete)
2. When operations conflict, we transform them against each other
3. The transformation adjusts positions to account for concurrent changes

## Project Structure

```
ot-collaborative-editor/
├── main.go      - WebSocket server and client management
├── ot.go        - OT core logic (Transform and Apply functions)
├── index.html   - Web client with collaborative editor
├── go.mod       - Go module dependencies
└── README.md    - This file
```

## How It Works

### 1. OT Core (ot.go)

The OT implementation has two main components:

#### Operations
```go
type Operation struct {
    Type   string  // "insert" or "delete"
    Pos    int     // Position in document
    Char   string  // Character to insert (for insert ops)
}
```

#### Transform Function
The `Transform()` function is the heart of OT. It takes two operations and transforms the first against the second:

```go
func Transform(op1, op2 Operation) Operation
```

**Example transformations:**

1. **Both insert at same position:**
   ```
   op1: Insert "A" at 5
   op2: Insert "B" at 5
   Result: op1 becomes Insert "A" at 6 (shifted right)
   ```

2. **Insert before another operation:**
   ```
   op1: Insert "A" at 10
   op2: Insert "B" at 5
   Result: op1 becomes Insert "A" at 11 (shifted right by 1)
   ```

3. **Delete before another operation:**
   ```
   op1: Insert "A" at 10
   op2: Delete at 5
   Result: op1 becomes Insert "A" at 9 (shifted left by 1)
   ```

### 2. WebSocket Server (main.go)

The server manages:
- **Hub**: Central coordinator for all clients and the document
- **Clients**: WebSocket connections for each user
- **Document**: Shared document state with version tracking

**Flow:**
1. Client connects → Receives current document state
2. Client makes edit → Sends operation to server
3. Server applies operation to document
4. Server broadcasts operation to all other clients
5. Other clients apply operation to their local copy

### 3. Web Client (index.html)

The client handles:
- **Change Detection**: Detects what the user typed/deleted
- **Operation Creation**: Converts changes into OT operations
- **Receiving Operations**: Applies remote operations while preserving cursor position
- **Visual Feedback**: Shows connection status and operation log

## Running the Project

### Prerequisites
- Go 1.22 or higher

### Steps

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Run the server:**
   ```bash
   go run .
   ```

3. **Open multiple browser windows:**
   - Open http://localhost:8080 in 2-3 browser windows/tabs
   - Each window represents a different user

4. **Start typing:**
   - Type in any window
   - Watch the text appear in real-time across all windows
   - Check the operation log to see OT in action

## Experimenting with OT

Try these scenarios to understand OT better:

### Scenario 1: Same Position Insertion
1. Open 2 windows
2. Place cursor at position 0 in both
3. Type "A" in window 1, then "B" in window 2
4. Notice how OT handles the conflict - both characters appear

### Scenario 2: Distributed Edits
1. Open 2 windows
2. Type at different positions simultaneously
3. Watch how operations don't interfere with each other

### Scenario 3: Delete While Inserting
1. Open 2 windows
2. In window 1: Insert text at position 5
3. In window 2: Delete at position 3
4. See how OT adjusts positions

## Understanding the Code

### Key Functions

**ot.go:**
- `Apply()` - Applies an operation to the document
- `Transform()` - Transforms one operation against another

**main.go:**
- `Hub.run()` - Main event loop handling clients and operations
- `Client.readPump()` - Reads operations from WebSocket
- `Client.writePump()` - Sends operations to WebSocket

**index.html (JavaScript):**
- `findInsertPosition()` - Detects where text was inserted
- `findDeletePosition()` - Detects where text was deleted
- `sendOperation()` - Sends operation to server
- `ws.onmessage` - Handles received operations

## Limitations & Simplifications

This is a educational implementation with some simplifications:

1. **No Undo/Redo**: Real OT systems need special handling for undo
2. **Character-level only**: Only handles single character insert/delete
3. **No persistence**: Document is lost when server restarts
4. **Basic conflict resolution**: Production systems use more sophisticated algorithms
5. **No network delay simulation**: Real systems must handle out-of-order messages

## Advanced OT Concepts (Not Implemented Here)

- **Intention Preservation**: Ensuring user intent is maintained
- **Convergence**: Guaranteeing all clients reach the same state
- **Tombstones**: Tracking deleted content for better transformation
- **Multiple operation types**: Formatting, styling, etc.
- **Compound operations**: Transforming multiple ops at once

## Further Learning

To dive deeper into OT:

1. **Google Wave OT Paper**: Original comprehensive OT research
2. **ShareJS**: Production-ready OT implementation
3. **Conflict-free Replicated Data Types (CRDTs)**: Alternative to OT

## Troubleshooting

**"Connection refused" error:**
- Make sure the server is running (`go run .`)
- Check that port 8080 is not in use

**Changes not syncing:**
- Open browser console (F12) to check for errors
- Verify WebSocket connection status in the UI

**Operations seem out of order:**
- This is normal - OT handles this
- Watch the operation log to see transformations

## License

This is an educational project - feel free to use and modify!
