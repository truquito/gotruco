#!/bin/bash

#SBATCH --job-name=gowalker-4p
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=2
#SBATCH --mem=1G
#SBATCH --time=5-00:00:00
#SBATCH --partition=normal
#SBATCH --qos=normal
#SBATCH --output=/clusteruy/home/juan.filevich/batches/out/%x.%j.out
#SBATCH --error=/clusteruy/home/juan.filevich/batches/out/%x.%j.out
#SBATCH --mail-type=BEGIN,END
#SBATCH --mail-user=juan.filevich@fing.edu.uy

# dump args
printf "starts: $(date)\n"
echo "args: ${@}"
echo "job id: $SLURM_JOB_ID"

# Directories
BASE_DIR=$HOME/Workspace/facu/gotruco
CHECKPOINT_DIR=$HOME/checkpoints/gowalker
SOURCE_FILE=cmd/walker/main.go
BINARY_NAME=walker

# Create checkpoint directory if it doesn't exist
mkdir -p $CHECKPOINT_DIR

# Checkpoint filename includes job ID to prevent conflicts
CHECKPOINT_FILE="$CHECKPOINT_DIR/game_checkpoint_$SLURM_JOB_ID.json"
# Use previous checkpoint if it exists and was specified
if [ -n "$1" ]; then
    PREV_CHECKPOINT="$1"
    if [ -f "$PREV_CHECKPOINT" ]; then
        echo "Using previous checkpoint: $PREV_CHECKPOINT"
        cp "$PREV_CHECKPOINT" "$CHECKPOINT_FILE"
    else
        echo "Warning: Previous checkpoint $PREV_CHECKPOINT not found, starting fresh"
    fi
fi

# Set time limit slightly shorter than SLURM limit
# 4.7 days = 406800 seconds (5 days = 432000 seconds)
TIME_LIMIT=406800

cd ${BASE_DIR}

# compile
go build -ldflags="-w -s" -o $BINARY_NAME $SOURCE_FILE

start_time=$(date +%s)

# Run the walker with checkpoint support and capture output
srun --export=ALL \
    ./$BINARY_NAME --checkpoint="$CHECKPOINT_FILE" --timelimit=$TIME_LIMIT
rc=$? # capture the walker’s exit status here

end_time=$(date +%s)
elapsed_time=$((end_time - start_time))

printf "done: $(date)\n"
printf "elapsed time: %d seconds\n" "$elapsed_time"

# Check if processing is complete or needs another job
if [ $rc -eq 0 ]; then
    echo "Game tree traversal complete! No need for another job."
else
    echo "Job time limit reached, checkpoint saved to: $CHECKPOINT_FILE"
    # Auto-resubmit (uncomment to enable)
    echo "Auto-resubmitting job..."
    sbatch $0 $CHECKPOINT_FILE
fi

# notice:
# an alternative to `if [ $rc -eq 0 ]; then` could be
# `if echo "$output" | grep -q "Processing complete"; then`
