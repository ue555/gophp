<?php
if (false) { echo "bad"; } elseif (true) {} elseif (1 / 0) { echo "bad"; } else { echo "bad"; }
if (false) {} else if (true) {} else if ($missing) { echo "bad"; }
echo "done\n";
