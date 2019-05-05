#!/bin/sh
cd "$(dirname "$0")"

# This script downloads the songs listed in the songs.txt file and convert them to wav
# it will always clean the dl directory

DL_DIR="../assets/dataset/downloads"
WAV_DIR="../assets/dataset/wav"

# Start by clean the download directory
rm -r $DL_DIR

# Create the wav directory
mkdir -p $WAV_DIR

# Download the songs in the dl directory
cat ../assets/songs.txt | xargs wget -q -P $DL_DIR

echo "Done downloading the songs !"

# Convert the mp3 files into wav files
for i in $DL_DIR/*.mp3; do
    # Remove path from name
    name="$(echo $i | sed "s#.*/##")"
    ffmpeg -i "$i" -acodec pcm_s16le -ar 44100 "$WAV_DIR/${name%.*}.wav"
done
