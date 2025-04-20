#!/bin/bash

# make sure the logs directory `--output` and `--output` exist before execution

# upload with:
# `rsync -avz --copy-links --progress -e 'ssh -p 10022' --exclude '__pycache__/' --exclude '*.out' --exclude '*.log' --exclude '.git/' $(pwd) 'cluster.uy:Workspace/facu'`

# when using standard args launch it as:
# `sbatch -J gowalker ~/Workspace/facu/gotruco/cmd/walker/run.slurm [<arg1> <arg2> <arg3>]`

# max time is `--time=5-00:00:00`

#SBATCH --job-name=gowalker-2p
#SBATCH --nodes=1 # nodes to use in total, per job or subjob
#SBATCH --ntasks=1 # when doing `srun` launch `ntaks` parallel jobs across `nodes`
#SBATCH --cpus-per-task=2
#SBATCH --mem=1G
#SBATCH --time=5-00:00:00
#SBATCH --partition=normal
#SBATCH --qos=normal
#SBATCH --output=/clusteruy/home/juan.filevich/batches/out/%x.%j.out
#SBATCH --error=/clusteruy/home/juan.filevich/batches/out/%x.%j.out
#SBATCH --mail-type=BEGIN,END,FAIL # options: NONE, BEGIN, END, FAIL, REQUEUE, ALL
#SBATCH --mail-user=juan.filevich@fing.edu.uy

# every (sub)job will execute this from here
# all (sub)job will share the same args ${@}

# dump args
printf "starts: $(date)\n"
echo "args: ${@}"

base=$HOME/Workspace/facu/gotruco # edit this!
cmd=cmd/walker/main.go

# srun will run this cmd `nodes` times in parallel
cd ${base}

# compile
go build -ldflags="-w -s" -o walker $cmd

start_time=$(date +%s)

srun --export=ALL \
    walker

end_time=$(date +%s)
elapsed_time=$((end_time - start_time))

printf "done: $(date)\n"
printf "elapsed time: %d seconds\n" "$elapsed_time"
