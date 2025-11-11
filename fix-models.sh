#!/bin/bash

# Move model files to the correct location
mv /home/josh/models/models/* /home/josh/models/

# Remove the empty subdirectory
rmdir /home/josh/models/models

echo "Model files moved to /home/josh/models/"