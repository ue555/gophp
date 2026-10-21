<?php
function ignoreWarning($level) { return null; }
set_error_handler('ignoreWarning');
echo $missing;
echo "done\n";
