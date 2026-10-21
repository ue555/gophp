<?php
echo "before\n";
function unused(): never { if (false) { return null; } throw new Exception('stop'); }
