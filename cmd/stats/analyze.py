import argparse
import json
import os
import matplotlib.pyplot as plt

def plot_stats(file_path):
    # 1. Load the file
    if not os.path.exists(file_path):
        print(f"Error: The file '{file_path}' was not found.")
        return

    with open(file_path, 'r') as f:
        try:
            stats = json.load(f)
        except json.JSONDecodeError:
            print(f"Error: Failed to decode JSON from '{file_path}'.")
            return

    # Helper to convert Go's string-keyed map back to sorted integer lists
    def get_sorted_xy(dist_map):
        if not dist_map:
            return [], []
        # Convert keys to int for numerical sorting
        sorted_keys = sorted([int(k) for k in dist_map.keys()])
        values = [dist_map[str(k)] for k in sorted_keys]
        return sorted_keys, values

    # 2. Extract Data
    dists = {
        "FloresDist": get_sorted_xy(stats.get('FloresDist', {})),
        "EnvidoDist": get_sorted_xy(stats.get('EnvidoDist', {})),
        "PoderDist":  get_sorted_xy(stats.get('PoderDist', {}))
    }

    # 3. Print Min/Max to stdout
    print(f"--- Statistics for Total: {stats.get('Total', 0)} ---")
    for name, (keys, _) in dists.items():
        if keys:
            print(f"{name}: Min Key = {min(keys)}, Max Key = {max(keys)}")
        else:
            print(f"{name}: (Map is empty)")

    # 4. Plotting (1 Row, 3 Columns)
    fig, axs = plt.subplots(1, 3, figsize=(18, 5))
    fig.suptitle(f"Distributions (Total: {stats.get('Total', 0)})", fontsize=16)
    
    colors = ['skyblue', 'salmon', 'lightgreen']
    
    for i, (name, (x, y)) in enumerate(dists.items()):
        axs[i].bar(x, y, color=colors[i])
        axs[i].set_title(name)
        axs[i].set_xlabel("Value")
        axs[i].set_ylabel("Frequency")
        axs[i].grid(axis='y', linestyle='--', alpha=0.7)

    plt.tight_layout(rect=[0, 0.03, 1, 0.95])
    plt.show()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Plot Go stats from a JSON file.')
    parser.add_argument('--json', type=str, required=True, help='Path to the JSON file')
    args = parser.parse_args()
    
    plot_stats(args.json)