
#!/bin/bash

cmd="/opt/orion/bin/orion /config/${ORION_NAME}/orion.conf"

echo ${cmd}

eval $cmd