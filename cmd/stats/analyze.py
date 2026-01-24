import argparse
import json
import os
import matplotlib.pyplot as plt

def plot_stats(file_path):
    if not os.path.exists(file_path):
        print(f"Error: The file '{file_path}' was not found.")
        return

    with open(file_path, 'r') as f:
        try:
            stats = json.load(f)
        except json.JSONDecodeError:
            print(f"Error: Failed to decode JSON from '{file_path}'.")
            return

    def get_sorted_xy(dist_map):
        if not dist_map:
            return [], []
        sorted_keys = sorted([int(k) for k in dist_map.keys()])
        values = [dist_map[str(k)] for k in sorted_keys]
        return sorted_keys, values

    # Extract Data
    f_x, f_y = get_sorted_xy(stats.get('FloresDist', {}))
    e_x, e_y = get_sorted_xy(stats.get('EnvidoDist', {}))
    p_x, p_y = get_sorted_xy(stats.get('PoderDist', {}))
    
    total = stats.get('Total', 0)
    dists_data = [
        ("FloresDist", f_x, f_y),
        ("EnvidoDist", e_x, e_y),
        ("PoderDist",  p_x, p_y)
    ]

    # --- Calculations & Console Output ---
    print(f"--- Statistics (Total Samples: {total}) ---")
    
    for name, x, y in dists_data:
        if x:
            print(f"{name}: Min Key = {min(x)}, Max Key = {max(x)}")
        else:
            print(f"{name}: (Map is empty)")

    # Probability of Flor calculation
    if total > 0:
        sum_flores = sum(f_y)
        prob_flor = sum_flores / total
        print(f"\nProbability of having a Flor: {prob_flor:.4f} ({prob_flor * 100:.2f}%)")
    else:
        print("\nProbability of Flor: N/A (Total is 0)")

    # --- Plotting (1 Row, 3 Columns) ---
    fig, axs = plt.subplots(1, 3, figsize=(18, 5))
    fig.suptitle(f"Distributions (Total: {total})", fontsize=16)
    
    colors = ['green', 'red', 'blue']
    
    for i, (name, x, y) in enumerate(dists_data):
        axs[i].bar(x, y, color=colors[i])
        axs[i].set_title(name)
        axs[i].set_xlabel("Key")
        axs[i].set_ylabel("Frequency")
        axs[i].grid(axis='y', linestyle='--', alpha=0.6)

    plt.tight_layout(rect=[0, 0.03, 1, 0.95])
    plt.show()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Plot Go stats from a JSON file.')
    parser.add_argument('--json', type=str, required=True, help='Path to the JSON file')
    args = parser.parse_args()
    
    plot_stats(args.json)