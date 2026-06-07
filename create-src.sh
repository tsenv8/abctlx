# SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
./abctlx sources create \
    --db="dellstore" \
    --host="localhost" \
    --name="abctlx-sourcedb" \
    --port="5432" \
    --pub="airbyte_publication" \
    --pw="postgres" \
    --rep="airbyte_slot" \
    --user="postgres"


#     Flags:
#       --db string        Database Name (default "postgres")
#   -h, --help             help for create
#       --host string      Database Host Name (default "localhost")
#       --name string      Source Name (default "sourcedb")
#       --port int         Connection Port (default 2499)
#       --pub string       Airbyte Publication Name (default "airbyte_publication")
#       --pw string        Database Password (default "1")
#       --rep string       Airbyte Replication Slot Name (default "airbyte_slot")
#       --schema strings   Database Schemas (default [public])
#       --user string      Database Username (default "postgres")
