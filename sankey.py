import csv
import plotly.graph_objects as go

# Path to the CSV file
csv_file_path = "connections.csv"

# Load data from CSV
connections = []
with open(csv_file_path, "r") as file:
    reader = csv.reader(file)
    first_row = next(reader)

    # Check if there is a header
    if "source" in first_row and "destination" in first_row and "protocol" in first_row:
        file.seek(0)
        reader = csv.DictReader(file)
        for row in reader:
            connections.append((row["source"], row["destination"], row["protocol"]))
    else:
        connections.append((first_row[0], first_row[1], first_row[2]))
        for row in reader:
            connections.append((row[0], row[1], row[2]))

# Extract nodes and flows for Sankey
sources = [src for src, dst, _ in connections]
destinations = [dst for src, dst, _ in connections]
protocols = [protocol for src, dst, protocol in connections]

# Create a unique list of all nodes (IPs)
nodes = list(set(sources + destinations))

# Map nodes to indices
node_indices = {node: idx for idx, node in enumerate(nodes)}

# Prepare data for Sankey diagram
sankey_sources = [node_indices[src] for src in sources]
sankey_targets = [node_indices[dst] for dst in destinations]
sankey_labels = nodes
sankey_values = [1] * len(sankey_sources)  # Assign a value of 1 for each connection

# Create Sankey diagram
fig = go.Figure(go.Sankey(
    node=dict(
        pad=15,
        thickness=20,
        line=dict(color="black", width=0.5),
        label=sankey_labels
    ),
    link=dict(
        source=sankey_sources,
        target=sankey_targets,
        value=sankey_values,
        label=protocols
    )
))

fig.update_layout(title_text="Sankey Diagram of Host Connections", font_size=10)
fig.show()
