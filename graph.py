import networkx as nx
import matplotlib.pyplot as plt
import csv

# Initialize directed graph
G = nx.DiGraph()

# Path to the CSV file
csv_file_path = "connections.csv"

# Define default column names
columns = ["source", "destination", "protocol"]

# Open the file and detect if there's a header
with open(csv_file_path, "r") as file:
    first_line = file.readline().strip().split(",")
    file.seek(0)  # Reset file pointer to the beginning

    if set(first_line) == set(columns):  # If the first line matches the columns
        reader = csv.DictReader(file)
    else:
        reader = csv.reader(file)
    
    # Read data from the CSV
    for row in reader:
        if isinstance(row, dict):  # If reading with DictReader
            src, dst, protocol = row["source"], row["destination"], row["protocol"]
        else:  # If reading with reader, use column indices
            src, dst, protocol = row[0], row[1], row[2]
        
        G.add_edge(src, dst, label=protocol)

# Draw the network diagram
plt.figure(figsize=(12, 8))
pos = nx.spring_layout(G, seed=42)  # Force-directed layout for readability
nx.draw(G, pos, with_labels=True, node_size=3000, node_color="lightblue", font_size=10, font_weight="bold", edge_color="gray", arrows=True)
edge_labels = nx.get_edge_attributes(G, 'label')
nx.draw_networkx_edge_labels(G, pos, edge_labels=edge_labels, font_size=8)

plt.title("Network Diagram of Host Connections")
plt.show()
