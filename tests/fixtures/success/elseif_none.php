<?php
if (0 > 1) { echo "bad"; } elseif (0 == 'word') { echo "bad"; } else if ('0') { echo "bad"; } else { echo "else|"; }
if (false) { echo 1 / 0; } elseif (null) { echo $missing; } else if (0.0) { echo "bad"; }
echo "done\n";
