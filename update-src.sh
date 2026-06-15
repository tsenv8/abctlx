# SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
./abctlx sources update \
    --target-source="asd" \
    --name="updated" \
    --port=8001 \ 

    # --db="dellstore" \
    # --host="localhost" \
    # --name="abctlx-sourcedb" \
    # --port="5432" \
    # --pub="airbyte_publication" \
    # --pw="postgres" \
    # --rep="airbyte_slot" \
    # --user="postgres"
